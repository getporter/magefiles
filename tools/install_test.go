package tools_test

import (
	"os"
	"path/filepath"
	"testing"

	"get.porter.sh/magefiles/tools"
	"github.com/carolynvs/magex/pkg"
	"github.com/carolynvs/magex/pkg/gopath"
	"github.com/carolynvs/magex/xplat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnsureKind(t *testing.T) {
	tmp, err := os.MkdirTemp("", "magefiles")
	require.NoError(t, err, "Error creating temp directory")
	defer os.RemoveAll(tmp)

	oldGoPath := os.Getenv("GOPATH")
	defer os.Setenv("GOPATH", oldGoPath)
	os.Setenv("GOPATH", tmp)

	oldPath := os.Getenv("PATH")
	defer os.Setenv("PATH", oldPath)
	os.Setenv("PATH", tmp)

	tools.EnsureKindAt(tools.DefaultKindVersion)
	xplat.PrependPath(gopath.GetGopathBin())

	require.FileExists(t, filepath.Join(tmp, "bin", "kind"+xplat.FileExt()))

	found, err := pkg.IsCommandAvailable("kind", "--version", tools.DefaultKindVersion)
	require.NoError(t, err, "IsCommandAvailable failed")
	assert.True(t, found, "kind was not available from its location in GOPATH/bin. PATH=%s", os.Getenv("PATH"))
}

func TestEnsureStaticCheck(t *testing.T) {
	tmp, err := os.MkdirTemp("", "magefiles")
	require.NoError(t, err, "Error creating temp directory")
	defer os.RemoveAll(tmp)

	oldGoPath := os.Getenv("GOPATH")
	defer os.Setenv("GOPATH", oldGoPath)
	os.Setenv("GOPATH", tmp)

	oldPath := os.Getenv("PATH")
	defer os.Setenv("PATH", oldPath)
	os.Setenv("PATH", tmp)

	tools.EnsureStaticCheck()
	xplat.PrependPath(gopath.GetGopathBin())

	require.FileExists(t, filepath.Join(tmp, "bin", "staticcheck"+xplat.FileExt()))

	found, err := pkg.IsCommandAvailable("staticcheck", "--version", tools.DefaultStaticCheckVersion)
	require.NoError(t, err, "IsCommandAvailable failed")
	assert.True(t, found, "staticcheck was not available from its location in GOPATH/bin. PATH=%s", os.Getenv("PATH"))
}

func TestEnsureGitHubClient(t *testing.T) {
	tmp, err := os.MkdirTemp("", "magefiles")
	require.NoError(t, err, "Error creating temp directory")
	defer os.RemoveAll(tmp)

	oldGoPath := os.Getenv("GOPATH")
	defer os.Setenv("GOPATH", oldGoPath)
	os.Setenv("GOPATH", tmp)

	oldPath := os.Getenv("PATH")
	defer os.Setenv("PATH", oldPath)
	os.Setenv("PATH", tmp)

	tools.EnsureGitHubClient()
	xplat.PrependPath(gopath.GetGopathBin())

	require.FileExists(t, filepath.Join(tmp, "bin", "gh"+xplat.FileExt()))

	found, err := pkg.IsCommandAvailable("gh", "--version", tools.DefaultGitHubClientVersion)
	require.NoError(t, err, "IsCommandAvailable failed")
	assert.True(t, found, "gh was not available from its location in GOPATH/bin. PATH=%s", os.Getenv("PATH"))
}
