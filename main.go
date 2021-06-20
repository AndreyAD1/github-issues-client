package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"gopl.io/ch4/github"
)

func getRepositoryIssues(ownerName string, repoName string) (
	[]github.Issue, error) {
	getRepoIssuesURL := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/issues",
		ownerName,
		repoName)

	var issues []github.Issue
	request, err := http.NewRequest("GET", getRepoIssuesURL, nil)
	if err != nil {
		return issues, err
	}
	request.Header.Set("Accept", "application/vnd.github.v3+json")
	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return issues, err
	}

	if response.StatusCode != http.StatusOK {
		response.Body.Close()
		err := fmt.Errorf("HTTP error: %s", response.Status)
		return issues, err
	}

	if err := json.NewDecoder(response.Body).Decode(&issues); err != nil {
		response.Body.Close()
		return issues, err
	}
	response.Body.Close()
	return issues, nil
}

func printRepositoryIssues(issues []github.Issue) {
	fmt.Printf("Total issue number: %d", len(issues))
	for i, issue := range issues {
		prettyIssue, err := json.MarshalIndent(issue, "", "\t")
		if err != nil {
			fmt.Printf("Can not prettify the issue number %d", issue.Number)
		}
		fmt.Printf("\nIssue no. %d\n", i)
		fmt.Println(string(prettyIssue))
	}
}

func main() {
	userName := flag.String("user", "", "a GitHub user name")
	password := flag.String("password", "", "a Github user password")
	ownerName := flag.String("owner", "", "owner of Github repository")
	repoName := flag.String("repo", "", "repository name")
	command := flag.String("command", "", "an action script should do")

	flag.Parse()

	argumentPerName := map[string]string{
		"user":     *userName,
		"password": *password,
		"owner":    *ownerName,
		"repo":     *repoName,
		"command":  *command,
	}
	var inputError string
	for name, value := range argumentPerName {
		if value == "" {
			inputError += "The script requires an argument " + name + "\n"
		}
	}
	if inputError != "" {
		fmt.Println(inputError)
		return
	}
	switch *command {
	case "repo-issues":
		issues, err := getRepositoryIssues(*ownerName, *repoName)
		if err != nil {
			fmt.Printf("ERROR: %s", err)
			return
		}
		printRepositoryIssues(issues)
	default:
		fmt.Println("Only 'repo-issues' command has been implemented yet")
	}
}
