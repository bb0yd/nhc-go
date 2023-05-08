package main

import (
	"fmt"
	"nhc-go/src/nvm_service"
	"os"
)

func main() {
	args := os.Args[1:]

	nodeVersion := ""
	if len(args) > 0 {
		nodeVersion = args[0]
		_, err := nvm_service.SetNodeVersion(nodeVersion)
		if err != nil {
			fmt.Println("Error setting Node version:", err)
			return
		}
	} else {
		nodeVersion, err := nvm_service.SetLatestNodeVersion()
		if err != nil {
			fmt.Println("Error setting latest Node version:", err)
			return
		}
		fmt.Printf("Successfully switched to Node version %s\n", nodeVersion)
	}
}
