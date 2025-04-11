package main

import (
	"fmt"

	jira "github.com/andygrunwald/go-jira"
)

// GetAllIssues will implement pagination of api and get all the issues.
// Jira API has limitation as to maxResults it can return at one time.
// You may have usecase where you need to get all the issues according to jql
// This is where this example comes in.
// GetAllIssues will implement pagination of api and get all the issues.
// Jira API has limitation as to maxResults it can return at one time.
// You may have usecase where you need to get all the issues according to jql
// This is where this example comes in.
// This function uses the new token-based pagination model introduced in Jira Cloud.
func GetAllIssues(client *jira.Client, searchString string) ([]jira.Issue, error) {
	var issues []jira.Issue

	// Create options with initial values
	opt := &jira.SearchOptions{
		MaxResults: 1000,             // Max results can go up to 1000
		Fields:     []string{"*all"}, // Request all navigable fields
	}

	// Loop until we have no more pages
	for {
		// Search with current options (token will be empty on first call)
		chunk, _, nextPageToken, err := client.Issue.Search(searchString, opt)
		if err != nil {
			return nil, err
		}

		// Initialize issues slice if needed
		if issues == nil {
			issues = make([]jira.Issue, 0)
		}

		// Add current chunk to results
		issues = append(issues, chunk...)

		// If no next page token, we're done
		if nextPageToken == "" {
			return issues, nil
		}

		// Set token for next page
		opt.NextPageToken = nextPageToken
	}
}

func main() {
	jiraClient, err := jira.NewClient(nil, "https://issues.apache.org/jira/")
	if err != nil {
		panic(err)
	}

	jql := "project = Mesos and type = Bug and Status NOT IN (Resolved)"
	fmt.Printf("Usecase: Running a JQL query '%s' with token-based pagination\n", jql)

	issues, err := GetAllIssues(jiraClient, jql)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Found %d issues matching the query\n", len(issues))

	// Print first few issues as an example
	for i, issue := range issues {
		if i >= 5 {
			fmt.Println("... and more")
			break
		}
		fmt.Printf("Issue %d: %s - %s\n", i+1, issue.Key, issue.Fields.Summary)
	}
}
