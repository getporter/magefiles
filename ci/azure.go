package ci

import (
	"fmt"
	"log"
	"os"

	"get.porter.sh/magefiles/tools"
	"github.com/uwu-tools/magex/pkg/gopath"

	"github.com/uwu-tools/magex/ci"
)

// ConfigureAgent sets up a CI worker agent with mage and ensures
// that GOPATH/bin is in PATH.
func ConfigureAgent() error {
	err := tools.EnsureMage()
	if err != nil {
		return err
	}

	// Instruct Azure DevOps to add GOPATH/bin to PATH
	gobin := gopath.GetGopathBin()
	err = os.MkdirAll(gobin, 07700)
	if err != nil {
		return fmt.Errorf("could not mkdir -p %s: %w", gobin, err)
	}

	p, _ := ci.DetectBuildProvider()
	log.Printf("Adding %s to the PATH\n", gobin)
	return p.PrependPath(gobin)
}
