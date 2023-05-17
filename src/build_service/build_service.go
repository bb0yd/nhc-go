package build_service

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

type BuildResult struct {
	InstallTime time.Duration
	Scripts     []ScriptResult
}

type ScriptResult struct {
	Name     string
	Duration time.Duration
	Success  bool
}

func BuildProject(projectsPath string) (bool, string, map[string]BuildResult) {
	directories, _ := ioutil.ReadDir(projectsPath)
	projectDirs := []string{}
	for _, dir := range directories {
		if dir.IsDir() {
			projectDirs = append(projectDirs, dir.Name())
		}
	}

	var wg sync.WaitGroup
	success := true
	results := make(map[string]chan BuildResult)

	for _, projectDir := range projectDirs {
		projectPath := filepath.Join(projectsPath, projectDir)
		if _, err := os.Stat(filepath.Join(projectPath, "package.json")); os.IsNotExist(err) {
			continue
		}

		results[projectDir] = make(chan BuildResult, 10)
		wg.Add(1)
		go func(projectDir string, projectPath string) {
			defer wg.Done()
			defer close(results[projectDir])

			installDependenciesStart := time.Now()
			cmd := exec.Command("npm", "i", "--prefix", projectPath)
			err := cmd.Run()
			installDependenciesDuration := time.Since(installDependenciesStart)

			result := BuildResult{
				InstallTime: installDependenciesDuration,
				Scripts:     []ScriptResult{},
			}

			if err != nil {
				success = false
				results[projectDir] <- result
				return
			}

			data, _ := ioutil.ReadFile(filepath.Join(projectPath, "package.json"))
			var pkg packageJSON
			json.Unmarshal(data, &pkg)

			buildScripts := []string{}
			for key := range pkg.Scripts {
				if strings.HasPrefix(key, "build") && key != "build:watch" {
					buildScripts = append(buildScripts, key)
				}
			}

			for _, buildScript := range buildScripts {
				buildStart := time.Now()
				cmd := exec.Command("npm", "run", buildScript, "--prefix", projectPath)
				err := cmd.Run()
				buildDuration := time.Since(buildStart)

				scriptResult := ScriptResult{
					Name:     buildScript,
					Duration: buildDuration,
					Success:  err == nil,
				}

				result.Scripts = append(result.Scripts, scriptResult)

				// Print formatted output
			fmt.Printf("Repository: %s\n", projectDir)
			fmt.Printf("Dependencies install time: %v\n", result.InstallTime)
			for _, script := range result.Scripts {
				status := "SUCCESS"
				if !script.Success {
					status = "FAILED"
				}
				fmt.Printf("Script: %s Duration: %v Result: %s\n", script.Name, script.Duration, status)
			}
			}

			results[projectDir] <- result
		}(projectDir, projectPath)
	}

	wg.Wait()

	buildResults := make(map[string]BuildResult)
	for project, resultChan := range results {
		for result := range resultChan {
			buildResults[project] = result
		}
	}

	// Print formatted output
	for repo, result := range buildResults {
		fmt.Printf("Repository: %s\n", repo)
		fmt.Printf("Dependencies install time: %v\n", result.InstallTime)
		for _, script := range result.Scripts {
			status := "SUCCESS"
			if !script.Success {
				status = "FAILED"
			}
			fmt.Printf("Script: %s Duration: %v Result: %s\n", script.Name, script.Duration, status)
		}
	}

	if success {
		return true, "All build scripts succeeded", buildResults
	} else {
		return false, "Some build scripts failed", buildResults
	}
}

