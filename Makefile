CURRENT_DIR=$(shell pwd)
DIST_DIR=${CURRENT_DIR}/dist

.PHONY: argocd-flux-plugin
argocd-flux-plugin:
	CGO_ENABLED=0 go build -v -i -o ${DIST_DIR}/argocd-flux-plugin ./cmd/argocd-flux-plugin
