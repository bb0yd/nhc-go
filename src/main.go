package main

import (
	"fmt"
	"os"
)

func main() {
	repoFile := "repositories.json"
	nodeVersion := "--lts"

	args := os.Args[1:]
	for i, arg := range args {
		if arg == "--repoFile" {
			repoFile = args[i+1]
		} else if arg == "--nodeVersion" {
			nodeVersion = args[i+1]
		}
	}

	nhc := NewNodeHealthCheck(repoFile, nodeVersion)
	err := nhc.Run()
	if err != nil {
		fmt.Printf("Error running Node Health Check: %v\n", err)
		os.Exit(1)
	}
}
