//go:build mage
// +build mage

package main

import (
	"get.porter.sh/magefiles/ci"
	"github.com/carolynvs/magex/mgx"
	"github.com/carolynvs/magex/shx"
)

var must = shx.CommandBuilder{StopOnError: true}

func ConfigureAgent() {
	mgx.Must(ci.ConfigureAgent())
}

func Build() {
	must.RunV("go", "build", "./...")
}

func Test() {
	must.RunV("go", "test", "./...")
}
