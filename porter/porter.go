package porter

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/carolynvs/magex/mgx"
	"github.com/carolynvs/magex/pkg/downloads"
	"github.com/carolynvs/magex/shx"
	"github.com/carolynvs/magex/xplat"
)

var (
	// DefaultPorterVersion is the default version of Porter that is installed when it's not present
	DefaultPorterVersion = "v1.0.0-alpha.19"
)

// Install the default version of porter
func EnsurePorter() {
	EnsurePorterAt(DefaultPorterVersion)
}

// Install the specified version of porter
func EnsurePorterAt(version string) {
	home := GetPorterHome()
	runtimesDir := filepath.Join(home, "runtimes")
	os.MkdirAll(runtimesDir, 0770)

	clientPath := filepath.Join(home, "porter"+xplat.FileExt())
	if _, err := os.Stat(clientPath); os.IsNotExist(err) {
		log.Println("Porter client not found at", clientPath)
		log.Println("Installing porter into", home)
		opts := downloads.DownloadOptions{
			UrlTemplate: "https://cdn.porter.sh/{{.VERSION}}/porter-{{.GOOS}}-{{.GOARCH}}{{.EXT}}",
			Name:        "porter",
			Version:     version,
			Ext:         xplat.FileExt(),
		}
		if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
			// we don't yet publish arm64 for porter
			opts.UrlTemplate = "https://cdn.porter.sh/{{.VERSION}}/porter-darwin-amd64"
		}
		err = downloads.Download(home, opts)
		mgx.Must(err)
	}

	runtimePath := filepath.Join(home, "runtimes", "porter-runtime")
	if _, err := os.Stat(runtimePath); os.IsNotExist(err) {
		log.Println("Porter runtime not found at", runtimePath)
		opts := downloads.DownloadOptions{
			UrlTemplate: "https://cdn.porter.sh/{{.VERSION}}/porter-linux-amd64",
			Name:        "porter-runtime",
			Version:     version,
		}
		err = downloads.Download(runtimesDir, opts)
		mgx.Must(err)
	}
}

type InstallMixinOptions struct {
	Name    string
	URL     string
	Feed    string
	Version string
}

// EnsureMixin installs the specified mixin.
func EnsureMixin(mixin InstallMixinOptions) error {
	home := GetPorterHome()
	mixinDir := filepath.Join(home, "mixins", mixin.Name)
	if _, err := os.Stat(mixinDir); err == nil {
		log.Println("Mixin already installed:", mixin.Name)
		return nil
	}

	log.Println("Installing mixin:", mixin.Name)
	if mixin.Version == "" {
		mixin.Version = DefaultMixinVersion
	}
	var source string
	if mixin.Feed != "" {
		source = "--feed-url=" + mixin.Feed
	} else {
		source = "--url=" + mixin.URL
	}

	porterPath := filepath.Join(home, "porter")
	return shx.Run(porterPath, "mixin", "install", mixin.Name, "--version", mixin.Version, source)
}
