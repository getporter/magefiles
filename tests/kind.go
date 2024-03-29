package tests

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"get.porter.sh/magefiles/docker"
	"get.porter.sh/magefiles/tools"
	"github.com/carolynvs/magex/mgx"
	"github.com/carolynvs/magex/pkg"
	"github.com/carolynvs/magex/shx"
	"github.com/magefile/mage/mg"
)

const (
	// Name of the KIND cluster used for testing
	DefaultKindClusterName = "porter"

	// Relative location of the KUBECONFIG for the test cluster
	Kubeconfig = "kind.config"
)

var (
	must = shx.CommandBuilder{StopOnError: true}

	//go:embed kind.config.yaml
	templateKindConfig string

	//go:embed local-registry.yaml
	templateLocalRegistry string
)

// Ensure that the test KIND cluster is up.
func EnsureTestCluster() {
	mg.Deps(EnsureKubectl)

	if !useCluster() {
		CreateTestCluster()
	}

	mgx.Must(docker.StartDockerRegistry())
}

// get the config of the current kind cluster, if available
func getClusterConfig() (kubeconfig string, ok bool) {
	contents, err := shx.OutputE("kind", "get", "kubeconfig", "--name", getKindClusterName())
	return contents, err == nil
}

// setup environment to use the current kind cluster, if available
func useCluster() bool {
	contents, ok := getClusterConfig()
	if ok {
		log.Println("Reusing existing kind cluster")

		userKubeConfig, _ := filepath.Abs(os.Getenv("KUBECONFIG"))
		currentKubeConfig := filepath.Join(pwd(), Kubeconfig)
		if userKubeConfig != currentKubeConfig {
			fmt.Printf("ATTENTION! You should set your KUBECONFIG to match the cluster used by this project\n\n\texport KUBECONFIG=%s\n\n", currentKubeConfig)
		}
		os.Setenv("KUBECONFIG", currentKubeConfig)

		if err := os.WriteFile(Kubeconfig, []byte(contents), 0660); err != nil {
			mgx.Must(fmt.Errorf("error writing %s: %w", Kubeconfig, err))
		}
		return true
	}

	return false
}

// Create a KIND cluster named porter.
func CreateTestCluster() {
	mg.Deps(tools.EnsureKind, docker.RestartDockerRegistry)

	// Determine host ip to populate kind config api server details
	// https://kind.sigs.k8s.io/docs/user/configuration/#api-server
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		mgx.Must(fmt.Errorf("could not get a list of network interfaces: %w", err))
	}

	var ipAddress string
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println("Current IP address : ", ipnet.IP.String())
				ipAddress = ipnet.IP.String()
				break
			}
		}
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		mgx.Must(fmt.Errorf("error retrieving the home directory of the current user %w", err))
	}
	kubeConfig := homeDir + "/.kube/config:" + filepath.Join(pwd(), Kubeconfig)

	os.Setenv("KUBECONFIG", kubeConfig)
	kindCfgTmpl, err := template.New("kind.config.yaml").Parse(templateKindConfig)
	if err != nil {
		mgx.Must(fmt.Errorf("error parsing EnsureKind config template: %w", err))
	}

	var kindCfgContents bytes.Buffer
	kindCfgData := struct {
		Address string
	}{
		Address: ipAddress,
	}
	if err = kindCfgTmpl.Execute(&kindCfgContents, kindCfgData); err != nil {
		mgx.Must(fmt.Errorf("could not render the kind.config template: %w", err))
	}
	if err = os.WriteFile("kind.config.yaml", kindCfgContents.Bytes(), 0660); err != nil {
		mgx.Must(fmt.Errorf("could not write kind config file: %w", err))
	}
	defer os.Remove("kind.config.yaml")

	mgx.Must(must.Command("kind", "create", "cluster", "--name", getKindClusterName(), "--config", "kind.config.yaml").
		Env("KIND_EXPERIMENTAL_DOCKER_NETWORK=" + docker.DefaultNetworkName).Run())

	// Document the local registry
	mgx.Must(kubectl("apply", "-f", "-").
		Stdin(strings.NewReader(templateLocalRegistry)).
		Run())
	mgx.Must(kubectl("config", "use-context", "kind-porter").Run())
}

// Delete the KIND cluster named porter.
func DeleteTestCluster() {
	mg.Deps(tools.EnsureKind)

	mgx.Must(must.RunE("kind", "delete", "cluster", "--name", getKindClusterName()))
}

func kubectl(args ...string) shx.PreparedCommand {
	kubeconfig := fmt.Sprintf("KUBECONFIG=%s", os.Getenv("KUBECONFIG"))
	return must.Command("kubectl", args...).Env(kubeconfig)
}

// Ensure kubectl is installed.
func EnsureKubectl() {
	if ok, _ := pkg.IsCommandAvailable("kubectl", "version", ""); ok {
		return
	}

	versionURL := "https://storage.googleapis.com/kubernetes-release/release/stable.txt"
	versionResp, err := http.Get(versionURL)
	if err != nil {
		mgx.Must(fmt.Errorf("unable to determine the latest version of kubectl: %w", err))
	}

	if versionResp.StatusCode > 299 {
		mgx.Must(fmt.Errorf("GET %s (%d): %s", versionURL, versionResp.StatusCode, versionResp.Status))
	}
	defer versionResp.Body.Close()

	kubectlVersion, err := io.ReadAll(versionResp.Body)
	if err != nil {
		mgx.Must(fmt.Errorf("error reading response from %s: %w", versionURL, err))
	}

	kindURL := "https://storage.googleapis.com/kubernetes-release/release/{{.VERSION}}/bin/{{.GOOS}}/{{.GOARCH}}/kubectl{{.EXT}}"
	mgx.Must(pkg.DownloadToGopathBin(kindURL, "kubectl", string(kubectlVersion)))
}

func pwd() string {
	wd, err := os.Getwd()
	if err != nil {
		mgx.Must(fmt.Errorf("pwd failed: %w", err))
	}
	return wd
}

func getKindClusterName() string {
	if name, ok := os.LookupEnv("KIND_NAME"); ok {
		return name
	}
	return DefaultKindClusterName
}
