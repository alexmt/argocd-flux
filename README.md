# argocd-flux

CLI tool which understands `.flux.yaml` format and generates K8S manifests in a Flux compatible way.

The `argocd-flux-plugin` is supposed to be configured as an Argo CD config management [plugin](https://argoproj.github.io/argo-cd/user-guide/config-management-plugins/):

```yaml
  configManagementPlugins: |
    - name: flux
      generate:
        command: [sh, -c]
        args: ["argocd-flux-plugin . --path $GIT_PATH"]
```

Below is a kustomize based deployment which deployes Argo CD with the configured `argocd-flux-plugin`. Live demo is available at: https://cd.apps.argoproj.io/applications/argocd-flux

`kustomization.yaml`:
```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

bases:
- github.com/argoproj/argo-cd//manifests/ha/cluster-install?ref=release-1.2

patchesStrategicMerge:
- argo-cd-cm.yaml
- argocd-repo-server-deploy.yaml

namespace: argocd
```

`argocd-repo-server-deploy.yaml`:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: argocd-repo-server
spec:
  template:
    spec:
      containers:
      - name: argocd-repo-server
        volumeMounts:
        - mountPath: /usr/local/bin/argocd-flux-plugin
          name: custom-tools
          subPath: argocd-flux-plugin
      volumes:
      - name: custom-tools
        emptyDir: {}
      initContainers:
      - name: download-tools
        image: alexmt/argocd-flux:latest
        command: [sh, -c]
        args:
        - cp /usr/local/bin/argocd-flux-plugin /custom-tools/argocd-flux-plugin
        volumeMounts:
        - mountPath: /custom-tools
          name: custom-tools

```

`argo-cd-cm.yaml`:
```yaml
apiVersion: v1
data:
  configManagementPlugins: |
    - name: flux
      generate:
        command: [sh, -c]
        args: ["argocd-flux-plugin . --path $GIT_PATH"]
kind: ConfigMap
metadata:
  name: argocd-cm
```
