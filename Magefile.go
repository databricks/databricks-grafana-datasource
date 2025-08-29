//go:build mage
// +build mage

package main

import (
	// mage:import
	build "github.com/grafana/grafana-plugin-sdk-go/build"
	"github.com/magefile/mage/mg"
)

// temporarily dropping LinuxARM target as Databricks SDK is not compatible
func BuildAllNoLinuxARM() {
	b := build.Build{}
	mg.Deps(b.Linux, b.Windows, b.Darwin, b.DarwinARM64, b.LinuxARM64)
}

// Default configures the default target.
var Default = BuildAllNoLinuxARM
