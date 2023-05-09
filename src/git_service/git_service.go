package git_service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type RepoInfo struct {
	Token        string   `json:"token"`
	Repositories []string `json:"repositories"`
}

type CloneResult struct {
	RepoURL  string
	Duration time.Duration
	Success  bool
	Size     string
}

func CloneReposFromFile(jsonFilePath string) ([]CloneResult, error) {
	repoInfos, err := readJSONFile(jsonFilePath)
	if err != nil {
		return nil, err
	}

	return cloneRepos(repoInfos)
}

func readJSONFile(jsonFilePath string) ([]RepoInfo, error) {
	data, err := ioutil.ReadFile(jsonFilePath)
	if err != nil {
		return nil, fmt.Errorf("unable to read JSON file: %v", err)
	}

	var repoInfos []RepoInfo
	err = json.Unmarshal(data, &repoInfos)
	if err != nil {
		return nil, fmt.Errorf("unable to parse JSON file: %v", err)
	}

	return repoInfos, nil
}

func cloneRepos(repoInfos []RepoInfo) ([]CloneResult, error) {
	clonedReposDir := "cloned_repos"
	err := os.RemoveAll(clonedReposDir)
	if err != nil {
		return nil, err
	}

	err = os.Mkdir(clonedReposDir, 0755)
	if err != nil {
		return nil, err
	}

	var results []CloneResult
	var wg sync.WaitGroup
	var concurrencyWg sync.WaitGroup
	resultsChan := make(chan CloneResult)

	concurrencyLimit := 5
	semaphore := make(chan struct{}, concurrencyLimit)

	for _, repoInfo := range repoInfos {
		for _, repoURL := range repoInfo.Repositories {
			wg.Add(1)
			concurrencyWg.Add(1)
			go func(token, repoURL string) {
				semaphore <- struct{}{}
				defer func() {
					<-semaphore
					concurrencyWg.Done()
				}()
				startTime := time.Now()
				err, repoPath := cloneRepo(token, repoURL, clonedReposDir)
				duration := time.Since(startTime)
				success := err == nil

				repoSize := "0"
				if success {
					repoSize, err = calculateRepoSize(repoPath)
				}

				resultsChan <- CloneResult{
					RepoURL:  repoURL,
					Duration: duration,
					Success:  success,
					Size:     repoSize,
				}

				if err != nil {
					fmt.Printf("Error cloning repository %s: %v\n", repoURL, err)
				}
				
				wg.Done()
			}(repoInfo.Token, repoURL)
		}
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	go func() {
		concurrencyWg.Wait()
		close(semaphore)
	}()

	for result := range resultsChan {
		fmt.Printf("Repository: %s, Duration: %v, Success: %v, Size: %s bytes\n", result.RepoURL, result.Duration, result.Success, result.Size)
		results = append(results, result)
	}

	return results, nil
}

func cloneRepo(token, repoURL, clonedReposDir string) (error, string) {
    gitCmd := exec.Command("git", "clone", repoURL)
    gitCmd.Dir = clonedReposDir
    gitCmd.Env = append(os.Environ(), fmt.Sprintf("GIT_ASKPASS=%s", token))
    err := gitCmd.Run()
    if err != nil {
        return fmt.Errorf("failed to clone repository %s: %v", repoURL, err), ""
    }

    repoPath := filepath.Join(clonedReposDir, filepath.Base(repoURL))
    if strings.HasSuffix(repoPath, ".git") {
        repoPath = repoPath[:len(repoPath)-4]
    }

    return nil, repoPath
}

func calculateRepoSize(repoPath string) (string, error) {
    var size int64
    err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() {
            size += info.Size()
        }
        return nil
    })

    if err != nil {
        return "", fmt.Errorf("failed to calculate repository size: %v", err)
    }

    units := []string{"B", "KB", "MB", "GB"}
    unitIndex := 0
    for size >= 1024 && unitIndex < len(units)-1 {
        size /= 1024
        unitIndex++
    }

    return fmt.Sprintf("%d %s", size, units[unitIndex]), nil
}
