//go:build mage
// +build mage

package main

import (
	"get.porter.sh/magefiles/ci"
	"get.porter.sh/magefiles/porter"
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

func Test() {
	// Set up the bin with porter installed, tests can use that to initialize their local test bin to save time re-downloading it
	mg.SerialDeps(porter.UseBinForPorterHome, porter.EnsurePorter)

	must.RunV("go", "test", "./...")
}
