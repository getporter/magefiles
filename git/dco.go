package git

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/carolynvs/magex/xplat"
)

//go:embed dco/prepare-commit-msg
var prepareCommitMsg string

// SetupDCO configures your git repository to automatically sign your commits
// to comply with our DCO
func SetupDCO() error {
	gotShell := xplat.DetectShell()
	if gotShell == "powershell" || gotShell == "cmd" {
		return fmt.Errorf("setupDCO must be run from a shell that supports bash but %s was detected", gotShell)
	}

	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("unable to determine the current working directory: %w", err)
	}
	repoRoot, err := FindRepositoryRoot(pwd)
	if err != nil {
		return err
	}
	hooksDir := filepath.Join(repoRoot, ".git/hooks/")
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		return fmt.Errorf("error ensuring that %s exists: %w", hooksDir, err)
	}
	hookPath := filepath.Join(hooksDir, "prepare-commit-msg")

	// Ask first before overwriting an existing hook
	_, err = os.Stat(hookPath)
	if err == nil {
		fmt.Printf("The DCO hook is already installed. Overwrite? [yN]: ")
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		if !strings.HasPrefix(strings.ToLower(response), "y") {
			fmt.Println("The DCO hook was not installed")
			return nil
		}
	} else if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("unable to read %s: %w", hookPath, err)
	}

	fmt.Println("Installing DCO git commit hook")
	return os.WriteFile(hookPath, []byte(prepareCommitMsg), 0755)
}
