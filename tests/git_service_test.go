package tests

import (
	"fmt"
	"os"
	"testing"

	"nhc-go/src/git_service"
)

const testJSONFile = "test_repos.json"

func TestCloneReposFromFile(t *testing.T) {
	// Prepare test JSON file
	prepareTestJSONFile()

	// Run CloneReposFromFile function
	results, err := git_service.CloneReposFromFile(testJSONFile)
	if err != nil {
		t.Errorf("Error during CloneReposFromFile: %v", err)
	}

	// Check results
	for _, result := range results {
		if !result.Success {
			t.Errorf("Clone failed for repository: %s", result.RepoURL)
		}
	}

	// Cleanup test JSON file and cloned repos
	cleanupTestFiles()
}

func prepareTestJSONFile() {
	testData := `[
		{
			"token": "your_token_here",
			"repositories": [
				"git@github.com:bb0yd/good.git",
                "git@github.com:bb0yd/bad.git"
			]
		}
	]`

	err := os.WriteFile(testJSONFile, []byte(testData), 0644)
	if err != nil {
		fmt.Printf("Error creating test JSON file: %v\n", err)
	}
}

func cleanupTestFiles() {
	err := os.Remove(testJSONFile)
	if err != nil {
		fmt.Printf("Error removing test JSON file: %v\n", err)
	}

	err = os.RemoveAll("cloned_repos")
	if err != nil {
		fmt.Printf("Error removing cloned repos directory: %v\n", err)
	}
}
