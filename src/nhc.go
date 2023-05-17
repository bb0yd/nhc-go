package main

import (
	"fmt"
	"nhc-go/src/build_service"
	"nhc-go/src/git_service"
	"nhc-go/src/nvm_service"
	"nhc-go/src/report_service"
)

type NodeHealthCheck struct {
	repoFile        string
	nodeVersion     string
}

func NewNodeHealthCheck(repoFile string, nodeVersion string) *NodeHealthCheck {
	return &NodeHealthCheck{
		repoFile:    repoFile,
		nodeVersion: nodeVersion,
	}
}

func (nhc *NodeHealthCheck) Run() error {
	err := nhc.setNodeVersion()
	if err != nil {
		return err
	}

	nhc.nodeVersion, err = nvm_service.GetNodeVersion()
	if err != nil {
		return err
	}

	cloneResults, err := nhc.cloneRepositories()
	if err != nil {
		return err
	}

	buildResults, err := nhc.buildProjects()
	if err != nil {
		return err
	}

	report_service.GenerateReport(nhc.nodeVersion, cloneResults, buildResults)

	return nil
}

func (nhc *NodeHealthCheck) setNodeVersion() error {
	var err error
	if nhc.nodeVersion == "--lts" {
		_, err = nvm_service.SetLatestNodeVersion()
	} else {
		_, err = nvm_service.SetNodeVersion(nhc.nodeVersion)
	}
	
	if err != nil {
		return fmt.Errorf("failed to set Node version: %v", err)
	}
	return nil
}

func (nhc *NodeHealthCheck) cloneRepositories() ([]git_service.CloneResult, error) {
	cloneResults, err := git_service.CloneReposFromFile(nhc.repoFile)
	if err != nil {
		return nil, fmt.Errorf("failed to clone repositories: %v", err)
	}
	return cloneResults, nil
}

func (nhc *NodeHealthCheck) buildProjects() (map[string]build_service.BuildResult, error) {
	projectsPath := "cloned_repos" // Provide the path to the projects directory
	_, _, buildResults := build_service.BuildProject(projectsPath)
	return buildResults, nil
}