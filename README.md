# kubectl-kubesec

This is a kubectl plugin for scanning Kubernetes deployments with [kubesec.io](https://kubesec.io)

### Install

Clone the repo in into your `GOPATH`, build the plugin and deploy it to `~/.kube/plugins/scan`:

```bash
git clone https://github.com/stefanprodan/kubectl-kubesec
cd kubectl-kubesec
make deploy
```

### Usage

Scan a deployment:

```bash
kubectl -n kube-system plugin scan kubernetes-dashboard
```

Result:

```
kubernetes-dashboard kubesec.io score 7
-----------------
Advise
1. containers[] .securityContext .runAsNonRoot == true
Force the running image to run as a non-root user to ensure least privilege
2. containers[] .securityContext .capabilities .drop
Reducing kernel capabilities available to a container limits its attack surface
3. containers[] .securityContext .readOnlyRootFilesystem == true
An immutable root filesystem can prevent malicious binaries being added to PATH and increase attack cost
4. containers[] .securityContext .runAsUser > 10000
Run as a high-UID user to avoid conflicts with the host's user table
5. containers[] .securityContext .capabilities .drop | index("ALL")
Drop all capabilities and add only those required to reduce syscall attack surface
```
