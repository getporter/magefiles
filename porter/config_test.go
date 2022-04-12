package porter

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/carolynvs/magex/xplat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPorterHome(t *testing.T) {
	t.Run("PORTER_HOME set", func(t *testing.T) {
		tmpPorterHome, err := ioutil.TempDir("", "magefiles")
		require.NoError(t, err)
		defer os.RemoveAll(tmpPorterHome)

		os.Setenv("PORTER_HOME", tmpPorterHome)
		defer os.Unsetenv("PORTER_HOME")

		gotHome := GetPorterHome()
		assert.Equal(t, tmpPorterHome, gotHome)
	})

	t.Run("Default to HOME/.porter", func(t *testing.T) {
		tmpUserHome, err := ioutil.TempDir("", "magefiles")
		require.NoError(t, err)
		defer os.RemoveAll(tmpUserHome)
		tmpPorterHome := filepath.Join(tmpUserHome, ".porter")
		err = os.Mkdir(tmpPorterHome, 0700)
		require.NoError(t, err)

		os.Setenv("HOME", tmpUserHome)
		defer os.Unsetenv("HOME")

		gotHome := GetPorterHome()
		assert.Equal(t, tmpPorterHome, gotHome)
	})

	t.Run("no home found", func(t *testing.T) {
		os.Unsetenv("HOME")

		defer func() {
			panicObj := recover()
			require.Contains(t, panicObj, "Could not find a Porter installation", "Expected a panic since there's no home for porter")
		}()

		GetPorterHome()
	})
}

func TestUseBinForPorterHome(t *testing.T) {
	defer os.RemoveAll("bin")

	UseBinForPorterHome()

	pwd, _ := os.Getwd()
	binDir := filepath.Join(pwd, "bin")
	assert.Equal(t, binDir, os.Getenv("PORTER_HOME"), "PORTER_HOME was not set correctly")

	assert.True(t, xplat.InPath(binDir), "expected the bin directory to be in the PATH")
}
