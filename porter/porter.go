package porter

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/carolynvs/magex/mgx"
	"github.com/carolynvs/magex/pkg"
	"github.com/carolynvs/magex/pkg/downloads"
	"github.com/carolynvs/magex/shx"
	"github.com/carolynvs/magex/xplat"
)

var (
	// DefaultPorterVersion is the default version of Porter that is installed when it's not present
	DefaultPorterVersion = "v1.0.0-alpha.19"
)

// Install the default version of porter, if porter isn't already installed
func EnsurePorter() {
	home := GetPorterHome()
	clientPath := filepath.Join(home, "porter"+xplat.FileExt())
	if ok, _ := pkg.IsCommandAvailable(clientPath, "--version", ""); !ok {
		EnsurePorterAt(DefaultPorterVersion)
	}
}

// Install the specified version of porter
func EnsurePorterAt(version string) {
	home := GetPorterHome()
	runtimesDir := filepath.Join(home, "runtimes")
	os.MkdirAll(runtimesDir, 0770)

	clientPath := filepath.Join(home, "porter"+xplat.FileExt())
	if clientFound, _ := pkg.IsCommandAvailable(clientPath, "--version", version); !clientFound {
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
		mgx.Must(downloads.Download(home, opts))
	}

	runtimePath := filepath.Join(home, "runtimes", "porter-runtime")
	if runtimeFound, _ := pkg.IsCommandAvailable(runtimePath, "--version", version); !runtimeFound {
		log.Println("Porter runtime not found at", runtimePath)
		opts := downloads.DownloadOptions{
			UrlTemplate: "https://cdn.porter.sh/{{.VERSION}}/porter-linux-amd64",
			Name:        "porter-runtime",
			Version:     version,
		}
		mgx.Must(downloads.Download(runtimesDir, opts))
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
