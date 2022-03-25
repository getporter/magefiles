//go:build mage
// +build mage

package main

import (
	"github.com/carolynvs/magex/shx"
)

var must = shx.CommandBuilder{StopOnError: true}

func Build() {
	must.RunV("go", "build", "./...")
}

func Test() {
	must.RunV("go", "test", "./...")
}
