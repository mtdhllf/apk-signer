package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func GetAPKPackage(apkPath string, aaptPath string) (string, error) {
	cmd := exec.Command(aaptPath, "dump", "badging", apkPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("aapt execution failed: %v", err)
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "package: name=") {
			parts := strings.Split(line, "'")
			if len(parts) >= 2 {
				return parts[1], nil
			}
		}
	}

	return "", fmt.Errorf("failed to get package name")
}

func IsPackageAllowed(pkg string) bool {
	if len(AllowedPackages) == 0 {
		return true // If no allowed packages are configured, all packages are allowed.
	}
	for _, allowed := range AllowedPackages {
		if pkg == allowed {
			return true
		}
	}
	return false
}
