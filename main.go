package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

var youtrackToGitLabUserMap = map[string]string{}

var youtrackToGitLabTypeMap = map[string]string{}

var youtrackToGitLabStateMap = map[string]string{}

func main() {
	// Replace these file paths with your input and output file paths
	youtrackCSVPath := `SOURCE PATH`
	gitlabCSVPath := `DEST PATH`

	// Open YouTrack CSV file
	youtrackFile, err := os.Open(youtrackCSVPath)
	if err != nil {
		fmt.Println("Error opening YouTrack CSV file:", err)
		return
	}
	defer youtrackFile.Close()

	// Open GitLab CSV file for writing
	gitlabFile, err := os.Create(gitlabCSVPath)
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

	fmt.Println("Conversion completed. GitLab import CSV file created at:", gitlabCSVPath)
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
	yt_project := strings.TrimSpace(youtrackRow[1])
	yt_summary := "[YT_ID " + yt_id + "]" + " [" + yt_project + "] " + strings.TrimSpace(youtrackRow[2])
	yt_reporter := strings.TrimSpace(youtrackRow[3])
	yt_type := strings.TrimSpace(youtrackRow[4])
	yt_state := strings.TrimSpace(youtrackRow[5])
	yt_assignee := strings.TrimSpace(youtrackRow[6])
	yt_description := "[Reported By: " + yt_reporter + "] " + strings.TrimSpace(youtrackRow[7])

	if yt_type == "" {
		return nil, false
	}

	// Map YouTrack Summary to GitLab title
	gitlabRow[0] = yt_summary

	if gitlabAssignee, ok := youtrackToGitLabUserMap[yt_assignee]; ok {
		yt_description = yt_description + "\n" + "/assign " + gitlabAssignee + "\n"
	}
	if gitlabType, ok := youtrackToGitLabTypeMap[yt_type]; ok {
		yt_description = yt_description + "\n" + "/label " + gitlabType + "\n"
	}
	if gitlabState, ok := youtrackToGitLabStateMap[yt_state]; ok {
		yt_description = yt_description + "\n" + "/label " + gitlabState + "\n"
	}

	gitlabRow[1] = `"` + yt_description + `"`

	return gitlabRow, true
}
