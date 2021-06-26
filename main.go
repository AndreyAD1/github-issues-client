package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"gopl.io/ch4/github"
)

const (
	gitHubAPIURL    = "https://api.github.com"
	createIssueHelp = "https://docs.github.com/en/rest/reference/issues#create-an-issue"
)

func getRepositoryIssues(ownerName string, repoName string) (
	[]github.Issue, error) {
	getRepoIssuesURL := fmt.Sprintf(
		"%s/repos/%s/%s/issues",
		gitHubAPIURL,
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
		fmt.Printf("\nIssue no. %d\n", i+1)
		fmt.Println(string(prettyIssue))
	}
}

func createIssue(
	username string,
	password string,
	ownerName string,
	repoName string,
	jsonIssueProperties string) (github.Issue, error) {
	createIssueURL := fmt.Sprintf(
		"%s/repos/%s/%s/issues",
		gitHubAPIURL,
		ownerName,
		repoName,
	)
	var issue github.Issue
	bodyReader := strings.NewReader(jsonIssueProperties)
	request, err := http.NewRequest("POST", createIssueURL, bodyReader)
	if err != nil {
		return issue, err
	}
	request.Header.Set("Accept", "application/vnd.github.v3+json")
	request.SetBasicAuth(username, password)
	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return issue, err
	}

	if response.StatusCode != http.StatusCreated {
		response.Body.Close()
		err := fmt.Errorf("HTTP error: %s", response.Status)
		return issue, err
	}

	if err := json.NewDecoder(response.Body).Decode(&issue); err != nil {
		response.Body.Close()
		return issue, err
	}
	response.Body.Close()
	return issue, nil
}

func updateIssue(
	username string,
	password string,
	ownerName string,
	repoName string,
	issueNumber uint64,
	jsonIssueProperties string) (github.Issue, error) {
	createIssueURL := fmt.Sprintf(
		"%s/repos/%s/%s/issues/%d",
		gitHubAPIURL,
		ownerName,
		repoName,
		issueNumber,
	)
	var issue github.Issue
	bodyReader := strings.NewReader(jsonIssueProperties)
	request, err := http.NewRequest("PATCH", createIssueURL, bodyReader)
	if err != nil {
		return issue, err
	}
	request.Header.Set("Accept", "application/vnd.github.v3+json")
	request.SetBasicAuth(username, password)
	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return issue, err
	}

	if response.StatusCode != http.StatusOK {
		response.Body.Close()
		err := fmt.Errorf("HTTP error: %s", response.Status)
		return issue, err
	}

	if err := json.NewDecoder(response.Body).Decode(&issue); err != nil {
		response.Body.Close()
		return issue, err
	}
	response.Body.Close()
	return issue, nil
}

func openFileInEditor(filename string) error {
	editorName := os.Getenv("EDITOR")
	if editorName == "" {
		editorName = "vim"
	}
	executable, err := exec.LookPath(editorName)
	if err != nil {
		return err
	}
	command := exec.Command(executable, filename)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}

func getEditorOutput() (string, error) {
	file, err := os.CreateTemp(os.TempDir(), "*")
	if err != nil {
		return "Can not create a temporary file", err
	}
	filename := file.Name()
	defer os.Remove(filename)
	
	if err = file.Close(); err != nil {
		errMsg := fmt.Sprintf(
			"A temporary file %s is already open",
			filename,
		)
		return errMsg, err
	}
	if err = openFileInEditor(filename); err != nil {
		return "Text editor error", err
	}
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Sprintf("Can not read the file %s", filename), err
	}
	
	return string(bytes), nil
}

func main() {
	userName := flag.String("user", "", "a GitHub user name")
	password := flag.String("password", "", "a Github user password")
	ownerName := flag.String("owner", "", "owner of Github repository")
	repoName := flag.String("repo", "", "repository name")
	_ = flag.NewFlagSet("repo-issues", flag.ExitOnError)
	createIssueCmd := flag.NewFlagSet("create-issue", flag.ExitOnError)
	updateIssueCmd := flag.NewFlagSet("update-issue", flag.ExitOnError)
	issueNumber := updateIssueCmd.Uint64("issue-number", 0, "An issue number")

	flag.Parse()

	argumentPerName := map[string]string{
		"user":     *userName,
		"password": *password,
		"owner":    *ownerName,
		"repo":     *repoName,
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
	if len(os.Args) < 10 {
		fmt.Println("expected 'repo-issues' or 'create-issue' subcommands")
		os.Exit(1)
	}
	switch os.Args[9] {
	case "repo-issues":
		issues, err := getRepositoryIssues(*ownerName, *repoName)
		if err != nil {
			fmt.Printf("ERROR: %s", err)
			return
		}
		printRepositoryIssues(issues)
	case "create-issue":
		createIssueCmd.Parse(os.Args[10:])
		issueProperties, err := getEditorOutput()
		if err != nil {
			fmt.Println(issueProperties, err)
			return
		}
		issue, err := createIssue(
			*userName,
			*password,
			*ownerName,
			*repoName,
			issueProperties,
		)
		if err != nil {
			fmt.Printf("ERROR: %s", err)
			return
		}
		prettyIssue, err := json.MarshalIndent(issue, "", "\t")
		if err != nil {
			fmt.Printf("Can not prettify the created issue %v\n", issue)
		}
		fmt.Printf("Created issue:\n%s", string(prettyIssue))
	case "update-issue":
		updateIssueCmd.Parse(os.Args[10:])
		if *issueNumber == 0 {
			fmt.Println("Add 'issue-number' argument")
			return
		}
		issueProperties, err := getEditorOutput()
		if err != nil {
			fmt.Println(issueProperties, err)
			return
		}
		issue, err := updateIssue(
			*userName,
			*password,
			*ownerName,
			*repoName,
			*issueNumber,
			issueProperties,
		)
		if err != nil {
			fmt.Printf("ERROR: %s", err)
			return
		}
		prettyIssue, err := json.MarshalIndent(issue, "", "\t")
		if err != nil {
			fmt.Printf("Can not prettify the created issue %v\n", issue)
		}
		fmt.Printf("Created issue:\n%s", string(prettyIssue))
	default:
		fmt.Printf("Command %s is not implemented", os.Args[9])
	}
}
