package nvm_service

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

func GetInstalledNodeVersions() ([]string, error) {
	output, err := executeNVMCommand("ls", "--no-colors")
	if err != nil {
		return nil, err
	}

	versions := strings.Split(output, "\n")
	var installedVersions []string
	for _, v := range versions {
		if strings.Contains(v, "v") && !strings.Contains(v, "->") {
			installedVersions = append(installedVersions, strings.TrimSpace(strings.Trim(v, " *")))
		}
	}
	return installedVersions, nil
}

// GetNodeVersion retrieves the current version of NodeJS installed via nvm.
func GetNodeVersion() (string, error) {
	output, err := executeNVMCommand("current")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}

func SetNodeVersion(version string) (string, error) {
    var command string

    if version == "--lts" {
        fmt.Println("Installing latest LTS version.")
        command = "install"
    } else {
        installedVersions, err := executeNVMCommand("ls", "--no-colors")
        if err != nil {
            return "", fmt.Errorf("failed to get installed Node versions: %v", err)
        }

        if strings.Contains(installedVersions, "v"+version) {
            fmt.Printf("v%s is already installed.\n", version)
            command = "use"
        } else {
            fmt.Printf("Installing %s.\n", version)
            command = "install"
        }
    }

    output, err := executeNVMCommand(command, version)
    if err != nil {
        return "", fmt.Errorf("failed to set Node version: %v", err)
    }

    versionRegex := regexp.MustCompile(`(v\d+\.\d+\.\d+)`)
    matchedVersion := versionRegex.FindString(output)

    if matchedVersion == "" {
        return "", fmt.Errorf("could not extract Node version from output")
    }

    fmt.Printf("Current version: %s\n", matchedVersion)
    return matchedVersion, nil
}

// SetLatestNodeVersion switches to the latest LTS version of NodeJS.
func SetLatestNodeVersion() (string, error) {
	output, err := executeNVMCommand("install", "--lts")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}

// executeNVMCommand executes a command with nvm in a Bash shell.
func executeNVMCommand(args ...string) (string, error) {
	cmd := exec.Command("bash", "-ic", fmt.Sprintf("source $NVM_DIR/nvm.sh; nvm %s", strings.Join(args, " ")))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to execute nvm command: %v, output: %s", err, string(output))
	}
	return string(output), nil
}


