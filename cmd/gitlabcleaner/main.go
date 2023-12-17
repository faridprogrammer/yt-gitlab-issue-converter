package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	gitLabURL      = "https://gitlab.com/api/v4"
	privateToken   = "your-token"
	projectID      = "your-project-id"
	issuesEndpoint = "/projects/" + projectID + "/issues"
)

func main() {
	// Get list of issues
	issues, err := getIssues()
	if err != nil {
		log.Fatal(err)
	}

	// Loop through issues and delete each one
	for _, issue := range issues {
		err := deleteIssue(issue.ID)
		if err != nil {
			log.Printf("Error deleting issue %d: %v", issue.ID, err)
		} else {
			fmt.Printf("Deleted issue %d\n", issue.ID)
		}
	}
}

// Issue represents a GitLab issue
type Issue struct {
	ID int `json:"id"`
}

// getIssues retrieves a list of issues from GitLab
func getIssues() ([]Issue, error) {
	url := gitLabURL + issuesEndpoint

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Private-Token", privateToken)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var issues []Issue
	if err := decodeJSON(resp.Body, &issues); err != nil {
		return nil, err
	}

	return issues, nil
}

// deleteIssue deletes an issue with the given ID
func deleteIssue(issueID int) error {
	url := fmt.Sprintf("%s%s/%d", gitLabURL, issuesEndpoint, issueID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Private-Token", privateToken)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete issue %d, status code: %d", issueID, resp.StatusCode)
	}

	return nil
}

func decodeJSON(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}
