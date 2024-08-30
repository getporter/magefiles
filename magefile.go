//go:build mage
// +build mage

package main

import (
	"get.porter.sh/magefiles/ci"
	"get.porter.sh/magefiles/git"
	"get.porter.sh/magefiles/porter"
	"get.porter.sh/magefiles/tools"
	"github.com/carolynvs/magex/mgx"
	"github.com/carolynvs/magex/shx"
	"github.com/magefile/mage/mg"
)

var must = shx.CommandBuilder{StopOnError: true}

func ConfigureAgent() {
	mgx.Must(ci.ConfigureAgent())
}

func Build() {
	must.RunV("go", "build", "./...")
}

func Vet() {
	must.RunV("go", "vet", "./...")
}

// Run golangci-lint on the project
func Lint() {
	mg.Deps(tools.EnsureGolangCILint)
	must.RunV("golangci-lint", "run", "./...")
}

func Test() {
	// Set up the bin with porter installed, tests can use that to initialize their local test bin to save time re-downloading it
	mg.SerialDeps(porter.UseBinForPorterHome, porter.EnsurePorter)

	must.RunV("go", "test", "./...")
}

// SetupDCO configures your git repository to automatically sign your commits
// to comply with our DCO
func SetupDCO() error {
	return git.SetupDCO()
}
