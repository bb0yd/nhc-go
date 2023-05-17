package report_service

import (
	"fmt"
	"io/ioutil"
	"nhc-go/src/build_service"
	"nhc-go/src/git_service"
	"time"
)

func GenerateReport(nodeVersion string, gitResults []git_service.CloneResult, buildResults map[string]build_service.BuildResult) error {
	report := "# Report\n\n"
	report += fmt.Sprintf("## Date and Time: %s\n\n", time.Now().Format(time.RFC1123))
	report += fmt.Sprintf("## Node Version: %s\n\n", nodeVersion)
	report += "## Downloaded Repositories\n\n"
	for _, result := range gitResults {
		if result.Success {
			report += fmt.Sprintf("- %s\n", result.RepoURL)
		}
	}
	report += "\n"

	report += "## Built Repositories\n\n"
	for project, result := range buildResults {
		report += fmt.Sprintf("- %s\n", project)
		report += fmt.Sprintf("  - Install Time: %v\n", result.InstallTime)
		for _, scriptResult := range result.Scripts {
			report += fmt.Sprintf("    - Script Name: %20.20s  Build Time: %20.20v  Success: %20.20v\n", scriptResult.Name, scriptResult.Duration, scriptResult.Success)
		}
		report += "\n"
	}
	report += "\n"

	report += "## Failed to Build\n\n"
	for project, result := range buildResults {
		failed := false
		for _, scriptResult := range result.Scripts {
			if !scriptResult.Success {
				report += fmt.Sprintf("- %s\n", project)
				failed = true
				break
			}
		}
		if failed {
			continue
		}
	}
	report += "\n"

	err := ioutil.WriteFile("report.md", []byte(report), 0644)
	if err != nil {
		return fmt.Errorf("failed to write report file: %v", err)
	}

	return nil
}
