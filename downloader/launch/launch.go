package launch

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func LaunchGameClient() error {
	exeName := "Project-Epoch.exe"
	_, err := os.Stat(exeName)
	if err != nil {
		return fmt.Errorf("unable to find %s: %w", exeName, err)
	}

	cmd := calculateLaunchCommand(exeName)
	if cmd == nil {
		return fmt.Errorf("failed to create launch command")
	}

	// Set working directory to current directory
	cmd.Dir, err = os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	fmt.Printf("Launching %s...\n", exeName)

	// Check if wine exists on Linux
	if runtime.GOOS != "windows" {
		_, err := exec.LookPath("wine")
		if err != nil {
			return fmt.Errorf("wine is not installed: %w", err)
		}
	}

	// Combine stdout and stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command synchronously
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run %s: %w", exeName, err)
	}

	fmt.Printf("Successfully ran %s\n", exeName)
	return nil
}
func calculateLaunchCommand(exeName string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.Command(exeName)
	} else {
		absPath, err := filepath.Abs(exeName)
		if err != nil {
			fmt.Printf("Failed to get absolute path for %s: %v\n", exeName, err)
			return nil
		}
		winePrefix := filepath.Join(filepath.Dir(absPath), ".wine")
		os.Setenv("WINEPREFIX", winePrefix)
		return exec.Command("wine", exeName)
	}
}
