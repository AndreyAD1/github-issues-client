package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopl.io/ch4/github"
)

func main() {
	getRepoIssuesURL := "https://api.github.com/repos/octocat/hello-world/issues"
	request, err := http.NewRequest("GET", getRepoIssuesURL, nil)
	if err != nil {
		fmt.Println(err)
	}
	request.Header.Set("Accept", "application/vnd.github.v3+json")
	response, err := http.DefaultClient.Do(request)

	if response.StatusCode != http.StatusOK {
		response.Body.Close()
		fmt.Printf("search query failed: %s", response.Status)
	}

	var result []github.Issue
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		response.Body.Close()
		fmt.Println(err)
	}
	response.Body.Close()
	fmt.Println(result)
}