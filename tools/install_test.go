package tools_test

import (
	"io/ioutil"
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
	tmp, err := ioutil.TempDir("", "magefiles")
	require.NoError(t, err, "Error creating temp directory")
	defer os.RemoveAll(tmp)

	oldGoPath := os.Getenv("GOPATH")
	defer os.Setenv("GOPATH", oldGoPath)

	os.Setenv("GOPATH", tmp)
	tools.EnsureKindAt(tools.DefaultKindVersion)
	xplat.PrependPath(gopath.GetGopathBin())

	require.FileExists(t, filepath.Join(tmp, "bin", "kind"+xplat.FileExt()))

	found, err := pkg.IsCommandAvailable("kind", tools.DefaultKindVersion, "--version")
	require.NoError(t, err, "IsCommandAvailable failed")
	assert.True(t, found, "kind was not available from its location in GOPATH/bin. PATH=%s", os.Getenv("PATH"))
}

func TestEnsureStaticCheck(t *testing.T) {
	tmp, err := ioutil.TempDir("", "magefiles")
	require.NoError(t, err, "Error creating temp directory")
	defer os.RemoveAll(tmp)

	oldGoPath := os.Getenv("GOPATH")
	defer os.Setenv("GOPATH", oldGoPath)

	os.Setenv("GOPATH", tmp)
	tools.EnsureStaticCheck()
	xplat.PrependPath(gopath.GetGopathBin())

	require.FileExists(t, filepath.Join(tmp, "bin", "staticcheck"+xplat.FileExt()))

	found, err := pkg.IsCommandAvailable("staticcheck", tools.DefaultStaticCheckVersion, "--version")
	require.NoError(t, err, "IsCommandAvailable failed")
	assert.True(t, found, "staticcheck was not available from its location in GOPATH/bin. PATH=%s", os.Getenv("PATH"))
}
