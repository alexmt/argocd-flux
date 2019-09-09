package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/spf13/cobra"
	"github.com/weaveworks/flux/cluster/kubernetes"
	"github.com/weaveworks/flux/cluster/kubernetes/resource"
	"github.com/weaveworks/flux/manifests"
)

type dummyNamespacer struct {
}

// Just return original resource namespace. Argo CD should take are of removing namespace for cluster level resources
func (ns *dummyNamespacer) EffectiveNamespace(manifest resource.KubeManifest, knownScopes kubernetes.ResourceScopes) (string, error) {
	return manifest.GetNamespace(), nil
}

func newCommand() *cobra.Command {
	var (
		subPaths string
	)
	var command = cobra.Command{
		Use: "argocd-flux REPO_PATH",
		RunE: func(c *cobra.Command, args []string) error {

			if len(args) != 1 {
				return errors.New("Repo path is not specified.")
			}
			rootPath, err := filepath.Abs(args[0])
			if err != nil {
				return err
			}
			var targetPaths []string
			for _, subPath := range strings.Split(subPaths, ",") {
				if subPath = strings.TrimSpace(subPath); subPath != "" {
					targetPaths = append(targetPaths, path.Join(rootPath, subPath))
				}
			}

			logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
			k8sManifests := kubernetes.NewManifests(&dummyNamespacer{}, logger)
			configAware, err := manifests.NewConfigAware(rootPath, targetPaths, k8sManifests)

			if err != nil {
				return err
			}

			resources, err := configAware.GetAllResourcesByID(context.Background())
			if err != nil {
				return err
			}

			for _, r := range resources {
				_, _ = os.Stdout.WriteString(fmt.Sprintf("---\n%s", string(r.Bytes())))
			}
			return nil
		},
	}
	command.Flags().StringVarP(&subPaths, "path", "p", "", "comma separated list of repository sub directories")
	return &command
}

func main() {
	if err := newCommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
