package tests

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

type packageJSON struct {
	Scripts map[string]string `json:"scripts"`
}

func BuildProject(projectsPath string) (bool, string, map[string]string) {
	directories, _ := ioutil.ReadDir(projectsPath)
	projectDirs := []string{}
	for _, dir := range directories {
		if dir.IsDir() {
			projectDirs = append(projectDirs, dir.Name())
		}
	}

	var wg sync.WaitGroup
	success := true
	results := make(map[string]chan string)

	for _, projectDir := range projectDirs {
		projectPath := filepath.Join(projectsPath, projectDir)
		if _, err := os.Stat(filepath.Join(projectPath, "package.json")); os.IsNotExist(err) {
			continue
		}

		results[projectDir] = make(chan string, 10)
		wg.Add(1)
		go func(projectDir string, projectPath string) {
			defer wg.Done()
			defer close(results[projectDir])

			installDependenciesStart := time.Now()
			cmd := exec.Command("npm", "i", "--prefix", projectPath)
			err := cmd.Run()
			installDependenciesDuration := time.Since(installDependenciesStart)

			if err != nil {
				success = false
				results[projectDir] <- fmt.Sprintf("Failed to install dependencies:\n%s", err.Error())
				return
			}

			data, _ := ioutil.ReadFile(filepath.Join(projectPath, "package.json"))
			var pkg packageJSON
			json.Unmarshal(data, &pkg)

			buildScripts := []string{}
			for key := range pkg.Scripts {
				if strings.HasPrefix(key, "build") {
					buildScripts = append(buildScripts, key)
				}
			}

			if len(buildScripts) == 0 {
				results[projectDir] <- fmt.Sprintf("Dependencies install time: %v\nNo build scripts found", installDependenciesDuration)
				return
			}

			for _, buildScript := range buildScripts {
				buildStart := time.Now()
				cmd := exec.Command("npm", "run", buildScript, "--prefix", projectPath)
				output, err := cmd.CombinedOutput()
				buildDuration := time.Since(buildStart)

				if err != nil {
					success = false
					errorLogFileName := fmt.Sprintf("%s_error.log", filepath.Join(projectPath, buildScript))
					ioutil.WriteFile(errorLogFileName, output, 0644)
					results[projectDir] <- fmt.Sprintf("Dependencies install time: %v\nBuild script: %s\nDuration: %v\nResult: FAILED", installDependenciesDuration, buildScript, buildDuration)
					return
				} else {
					results[projectDir] <- fmt.Sprintf("Dependencies install time: %v\nBuild script: %s\nDuration: %v\nResult: SUCCESS", installDependenciesDuration, buildScript, buildDuration)
				}
			}
		}(projectDir, projectPath)
	}

	wg.Wait()

	buildResults := make(map[string]string)
	for project, resultChan := range results {
		for result := range resultChan {
			buildResults[project] = result
		}
	}

	if success {
		return true, "All build scripts succeeded", buildResults
	} else {
		return false, "Some build scripts failed", buildResults
	}
}
