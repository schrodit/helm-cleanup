[![build](https://github.com/schrodit/helm-cleanup/actions/workflows/build.yml/badge.svg)](https://github.com/schrodit/helm-cleanup/actions/workflows/build.yml)
[![lint](https://github.com/schrodit/helm-cleanup/actions/workflows/lint.yml/badge.svg)](https://github.com/schrodit/helm-cleanup/actions/workflows/lint.yml)

# helm-cleanup

## About

Identifies and cleans up Kubernetes resources that have been managed by Helm.

### Features

- Scans your cluster for resources managed by Helm that are no longer tracked by any release.
- Supports dry-run mode to preview actions before making changes.
- Interactive confirmation for resource deletion.
- Works with custom kubeconfig and context settings.
- Integrates with Helm's plugin system.


## Installation

Based on the version in plugin.yaml, release binary will be downloaded from GitHub:

```console
$ helm plugin install https://github.com/schrodit/helm-cleanup
Downloading and installing helm-cleanup v0.1.0 ...
https://github.com/helm/helm-cleanup/releases/download/v0.1.0/helm-cleanup_0.1.0_darwin_amd64.tar.gz
Installed plugin: cleanup
```

## Usage

```console
$ helm cleanup --help                                                                                                                                                                                                                                                                                  ☸dc5-sso|codesphere
Identifies and cleans up resources that have been managed by Helm.

Usage:
  cleanup [flags]

Flags:
      --dry-run               simulate a command (default true)
  -h, --help                  help for cleanup
      --kube-context string   name of the kubeconfig context to use
      --kubeconfig string     path to the kubeconfig file
      --namespace string      namespace scope of the release
      --yes                   Do not prompt on every deleted resource
```

:warning: By default this plugins runs in `dry-run` mode to not delete resources. Make sure to set `--dry-run=false` to cleanup leaked resources.

Example output:

```console
$ helm cleanup                                                                                                                                                                                                                                                                                         ☸dc5-sso|codesphere
NOTE: This is in dry-run mode, the following actions will not be executed.
Run without --dry-run to take the actions described below:
                                                                                                                                                                                                                                                                     | (0/76, 0 it/hr) [0s:0s]I0808 13:58:02.895047  239689 warnings.go:110] "Warning: v1 ComponentStatus is deprecated in v1.19+"
Fetch K8s resources 100% |███████████████| (76/76, 2 it/s)
Found 2 leaked resources
+-------------------------+--------------------------+-----------+---------------------------------+------------------------------+
| API VERSION             | KIND                     | NAMESPACE | NAME                            | HELM RELEASE                 |
+-------------------------+--------------------------+-----------+---------------------------------+------------------------------+
| apiextensions.k8s.io/v1 | CustomResourceDefinition |           | dnsentries.dns.gardener.cloud   | dns-controller/garden-system |
| apiextensions.k8s.io/v1 | CustomResourceDefinition |           | dnsproviders.dns.gardener.cloud | dns-controller/garden-system |
+-------------------------+--------------------------+-----------+---------------------------------+------------------------------+
```
