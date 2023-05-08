package tests

import (
	"fmt"
	"nhc-go/src/nvm_service"
	"testing"
)

func TestGetNodeVersion(t *testing.T) {
	_, err := nvm_service.GetNodeVersion()
	if err != nil {
		t.Errorf("GetNodeVersion failed with error: %v", err)
	}
}

func TestSetNodeVersion(t *testing.T) {
	testVersion := "v18.16.0"
	fmt.Println("Test version:", testVersion)

	currentVersion, err := nvm_service.SetNodeVersion(testVersion)
	if err != nil {
		t.Fatalf("SetNodeVersion failed with error: %v", err)
	}

	if currentVersion != testVersion {
		t.Fatalf("SetNodeVersion failed to set the correct version. Expected: %s, Got: %s", testVersion, currentVersion)
	}
}

func TestSetLatestNodeVersion(t *testing.T) {
	_, err := nvm_service.SetLatestNodeVersion()
	if err != nil {
		t.Errorf("SetLatestNodeVersion failed with error: %v", err)
	}
}
