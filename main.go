package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/xanzy/go-gitlab"
)

func main() {
	mrs, err := getMergeRequests()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v", mrs)
}

func getMergeRequests() (mrs []*gitlab.MergeRequest, e error) {

	projectId := os.Args[1]
	gitlabInstance := os.Args[2]
	authToken := os.Args[3]

	branchName, err := getCurrentBranch()
	if err != nil {
		return nil, err
	}

	var apiCustUrl = fmt.Sprintf(gitlabInstance + "/api/v4")

	git, err := gitlab.NewClient(authToken, gitlab.WithBaseURL(apiCustUrl))

	if err != nil {
		return nil, fmt.Errorf("Failed to create client: %v", err)
	}

	options := gitlab.ListProjectMergeRequestsOptions{
		Scope:        gitlab.String("all"),
		State:        gitlab.String("opened"),
		SourceBranch: &branchName,
	}

	mergeRequests, _, err := git.MergeRequests.ListProjectMergeRequests(projectId, &options)
	if err != nil {
		return nil, fmt.Errorf("Failed to list merge requests: %w", err)
	}

	if len(mergeRequests) == 0 {
		return nil, errors.New("No merge requests found")
	}

	return mergeRequests, nil
}

/* Gets the current branch */
func getCurrentBranch() (res string, e error) {
	gitCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")

	output, err := gitCmd.Output()
	if err != nil {
		return "", fmt.Errorf("Error running git rev-parse: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}
