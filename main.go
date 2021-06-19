package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	"gopl.io/ch4/github"
)

func main() {
	userName := flag.String("user", "", "a GitHub user name")
	password := flag.String("password", "", "a Github user password")
	ownerName := flag.String("owner", "", "owner of Github repository")
	repoName := flag.String("repo", "", "repository name")
	command := flag.String("command", "", "an action name script should do")

	flag.Parse()

	argumentPerName := map[string]string {
		"user": *userName,
		"password": *password,
		"owner": *ownerName,
		"repo": *repoName,
		"command": *command,
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

	getRepoIssuesURL := os.Args[1]
	request, err := http.NewRequest("GET", getRepoIssuesURL, nil)
	if err != nil {
		fmt.Println(err)
	}
	request.Header.Set("Accept", "application/vnd.github.v3+json")
	response, err := http.DefaultClient.Do(request)

	if err != nil {
		fmt.Println("Can not download issues")
		return
	}

	if response.StatusCode != http.StatusOK {
		response.Body.Close()
		fmt.Printf("search query failed: %s", response.Status)
	}

	var issues []github.Issue
	if err := json.NewDecoder(response.Body).Decode(&issues); err != nil {
		response.Body.Close()
		fmt.Println(err)
	}
	response.Body.Close()
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
