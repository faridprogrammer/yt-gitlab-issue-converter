package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	SrcPath               string
	DestPath              string
	YoutrackToGitlabUser  map[string]string
	YoutrackToGitlabType  map[string]string
	YoutrackToGitlabState map[string]string
}

var configObj Config

func main() {
	configPath := flag.String("c", "config.json", "Path to config file")

	flag.Parse()

	configContent, err := os.ReadFile(*configPath)

	if err != nil {
		fmt.Println("Error in opening config file:", err)
		return
	}

	err = json.Unmarshal(configContent, &configObj)

	if err != nil {
		fmt.Println("Error in reading config file:", err)
		return
	}

	// Open YouTrack CSV file
	youtrackFile, err := os.Open(configObj.SrcPath)
	if err != nil {
		fmt.Println("Error opening YouTrack CSV file:", err)
		return
	}
	defer youtrackFile.Close()

	// Open GitLab CSV file for writing
	gitlabFile, err := os.Create(configObj.DestPath)
	if err != nil {
		fmt.Println("Error creating GitLab CSV file:", err)
		return
	}
	defer gitlabFile.Close()

	// Create CSV reader and writer
	youtrackReader := csv.NewReader(youtrackFile)
	gitlabWriter := csv.NewWriter(gitlabFile)

	if err != nil {
		fmt.Println("Error reading YouTrack CSV header:", err)
		return
	}

	// Map YouTrack fields to GitLab fields
	gitlabHeader := getGitlabHeader()

	// Write GitLab header
	err = gitlabWriter.Write(gitlabHeader)
	if err != nil {
		fmt.Println("Error writing GitLab CSV header:", err)
		return
	}

	// Process each row
	firstRowPassed := false

	for {
		row, err := youtrackReader.Read()
		if err != nil {
			break // End of file
		}

		if !firstRowPassed {
			firstRowPassed = true
			continue
		}

		// Map YouTrack fields to GitLab fields, including assigning GitLab users
		gitlabRow, success := mapYouTrackToGitLab(row)
		if !success {
			continue
		}

		// Write GitLab row
		err = gitlabWriter.Write(gitlabRow)
		if err != nil {
			fmt.Println("Error writing GitLab CSV row:", err)
			return
		}
	}

	// Flush the writer to ensure all data is written to the file
	gitlabWriter.Flush()

	fmt.Println("Conversion completed. GitLab import CSV file created at:", configObj.DestPath)
}

func getGitlabHeader() []string {

	gitlabRow := make([]string, 2)

	gitlabRow[0] = "title"

	gitlabRow[1] = "description"

	return gitlabRow

}

func mapYouTrackToGitLab(youtrackRow []string) ([]string, bool) {

	gitlabRow := make([]string, len(youtrackRow))

	yt_id := strings.TrimSpace(youtrackRow[0])
	// yt_project := strings.TrimSpace(youtrackRow[1])
	yt_summary := yt_id + " " + strings.TrimSpace(youtrackRow[2])
	yt_reporter := strings.TrimSpace(youtrackRow[3])
	yt_type := strings.TrimSpace(youtrackRow[4])
	yt_state := strings.TrimSpace(youtrackRow[5])
	yt_assignee := strings.TrimSpace(youtrackRow[6])
	yt_description := "[Reported By: " + yt_reporter + "] \n" + strings.TrimSpace(youtrackRow[7])

	if yt_type == "" {
		return nil, false
	}

	// Map YouTrack Summary to GitLab title
	gitlabRow[0] = yt_summary

	if gitlabAssignee, ok := configObj.YoutrackToGitlabUser[yt_assignee]; ok {
		yt_description = yt_description + "\n" + "/assign " + gitlabAssignee + "\n"
	}
	if gitlabType, ok := configObj.YoutrackToGitlabType[yt_type]; ok {
		yt_description = yt_description + "\n" + "/label " + gitlabType + "\n"
	}
	if gitlabState, ok := configObj.YoutrackToGitlabState[yt_state]; ok {
		yt_description = yt_description + "\n" + "/label " + gitlabState + "\n"
	}

	gitlabRow[1] = `"` + yt_description + `"`

	return gitlabRow, true
}
