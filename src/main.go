package main

import (
	"fmt"
	"nhc-go/src/git_service"
	"os"
	"time"
)

func main() {
	start := time.Now()
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <path_to_json_file>")
		return
	}

	jsonFilePath := os.Args[1]
	_, err := git_service.CloneReposFromFile(jsonFilePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// fmt.Println("Clone Results:")
	// for _, result := range results {
	// 	fmt.Printf("Repository: %s, Duration: %v, Success: %v\n", result.RepoURL, result.Duration, result.Success)
	// }

	elapsed := time.Since(start)
	fmt.Printf("Program took %s to run\n", elapsed)
}


// NVM DRIVER
// package main

// import (
// 	"fmt"
// 	"nhc-go/src/nvm_service"
// 	"os"
// )

// func main() {
// 	args := os.Args[1:]

// 	nodeVersion := ""
// 	if len(args) > 0 {
// 		nodeVersion = args[0]
// 		_, err := nvm_service.SetNodeVersion(nodeVersion)
// 		if err != nil {
// 			fmt.Println("Error setting Node version:", err)
// 			return
// 		}
// 	} else {
// 		nodeVersion, err := nvm_service.SetLatestNodeVersion()
// 		if err != nil {
// 			fmt.Println("Error setting latest Node version:", err)
// 			return
// 		}
// 		fmt.Printf("Successfully switched to Node version %s\n", nodeVersion)
// 	}
// }
