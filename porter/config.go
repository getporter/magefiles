package porter

import (
	"os"
	"path/filepath"

	"github.com/carolynvs/magex/xplat"
)

var (
	// DefaultMixinVersion is the default version of mixins installed when it's not present
	DefaultMixinVersion = "canary"
)

// GetPorterHome determines the current PORTER_HOME directory
func GetPorterHome() string {
	porterHome := os.Getenv("PORTER_HOME")
	if porterHome == "" {
		home, _ := os.UserHomeDir()
		porterHome = filepath.Join(home, ".porter")
	}
	if _, err := os.Stat(porterHome); err != nil {
		panic("Could not find a Porter installation. Make sure that Porter is installed and set PORTER_HOME if you are using a non-standard installation path")
	}
	return porterHome
}

// UseBinForPorterHome sets the bin/ directory to be PORTER_HOME
func UseBinForPorterHome() {
	// use bin as PORTER_HOME
	wd, _ := os.Getwd()
	home := filepath.Join(wd, "bin")
	os.Mkdir(home, 0770)
	UsePorterHome(home)
}

// UsePorterHome sets the specified directory to be PORTER_HOME
func UsePorterHome(home string) {
	os.Setenv("PORTER_HOME", home)

	// Add PORTER_HOME to the PATH
	xplat.EnsureInPath(home)
}
