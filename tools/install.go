package tools

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/carolynvs/magex/mgx"
	"github.com/carolynvs/magex/pkg"
	"github.com/carolynvs/magex/pkg/archive"
	"github.com/carolynvs/magex/pkg/downloads"
	"github.com/carolynvs/magex/shx"
)

var (
	must = shx.CommandBuilder{StopOnError: true}

	// DefaultKindVersion is the default version of KinD that is installed when it's not present
	DefaultKindVersion = "v0.12.0"

	// DefaultStaticCheckVersion is the default version of StaticCheck that is installed when it's not present
	DefaultStaticCheckVersion = "2022.1.2"
)

// Fail if the go version doesn't match the specified constraint
// Examples: >=1.16
func EnforceGoVersion(constraint string) {
	log.Printf("Checking go version against constraint %s...", constraint)

	value := strings.TrimPrefix(runtime.Version(), "go")
	version, err := semver.NewVersion(value)
	if err != nil {
		mgx.Must(fmt.Errorf("could not parse go version: '%s': %w", value, err))
	}
	versionCheck, err := semver.NewConstraint(constraint)
	if err != nil {
		mgx.Must(fmt.Errorf("invalid semver constraint: '%s': %w", constraint, err))
	}

	ok, _ := versionCheck.Validate(version)
	if !ok {
		mgx.Must(fmt.Errorf("your version of Go, %s, does not meet the requirement %s", version, versionCheck))
	}
}

// Install mage
func EnsureMage() error {
	return pkg.EnsureMage("")
}

// Install gh
func EnsureGitHubClient() {
	if ok, _ := pkg.IsCommandAvailable("gh", ""); ok {
		return
	}

	// gh cli unfortunately uses a different archive schema depending on the OS
	target := "gh_{{.VERSION}}_{{.GOOS}}_{{.GOARCH}}/bin/gh{{.EXT}}"
	if runtime.GOOS == "windows" {
		target = "bin/gh.exe"
	}

	opts := archive.DownloadArchiveOptions{
		DownloadOptions: downloads.DownloadOptions{
			UrlTemplate: "https://github.com/cli/cli/releases/download/v{{.VERSION}}/gh_{{.VERSION}}_{{.GOOS}}_{{.GOARCH}}{{.EXT}}",
			Name:        "gh",
			Version:     "1.8.1",
			OsReplacement: map[string]string{
				"darwin": "macOS",
			},
		},
		ArchiveExtensions: map[string]string{
			"linux":   ".tar.gz",
			"darwin":  ".tar.gz",
			"windows": ".zip",
		},
		TargetFileTemplate: target,
	}

	if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
		// github doesn't publish arm64 binaries for mac, so fallback to the amd64 binary
		opts.UrlTemplate = "https://github.com/cli/cli/releases/download/v{{.VERSION}}/gh_{{.VERSION}}_{{.GOOS}}_amd64{{.EXT}}"
		opts.TargetFileTemplate = "gh_{{.VERSION}}_{{.GOOS}}_amd64/bin/gh{{.EXT}}"
	}

	err := archive.DownloadToGopathBin(opts)
	mgx.Must(err)
}

// Install kind
func EnsureKind() {
	EnsureKindAt(DefaultKindVersion)
}

// Install kind at the specified version
func EnsureKindAt(version string) {
	if ok, _ := pkg.IsCommandAvailable("kind", version); ok {
		return
	}

	kindURL := "https://github.com/kubernetes-sigs/kind/releases/download/{{.VERSION}}/kind-{{.GOOS}}-{{.GOARCH}}"
	mgx.Must(pkg.DownloadToGopathBin(kindURL, "kind", version))
}

// Install Staticcheck
func EnsureStaticCheck() {
	EnsureStaticCheckAt(DefaultStaticCheckVersion)
}

// Install Staticcheck at the specified version
func EnsureStaticCheckAt(version string) {
	if ok, _ := pkg.IsCommandAvailable("staticcheck", version); ok {
		return
	}

	opts := archive.DownloadArchiveOptions{
		DownloadOptions: downloads.DownloadOptions{
			UrlTemplate: "https://github.com/dominikh/go-tools/releases/download/{{.VERSION}}/staticcheck_{{.GOOS}}_{{.GOARCH}}.tar.gz",
			Name:        "staticcheck",
			Version:     version,
		},
		ArchiveExtensions: map[string]string{
			"linux":   ".tar.gz",
			"darwin":  ".tar.gz",
			"windows": ".tar.gz",
		},
		TargetFileTemplate: "staticcheck/staticcheck{{.EXT}}",
	}
	err := archive.DownloadToGopathBin(opts)
	mgx.Must(err)
}
