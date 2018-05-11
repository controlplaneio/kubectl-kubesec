# kubectl-kubesec

This is a kubectl plugin for scanning Kubernetes deployments with [kubesec.io](https://kubesec.io)

### Install

Download and extract the scan plugin to `~/.kube/plugins/scan`:

```bash
mkdir -p ~/.kube/plugins/scan && \
curl -sL https://github.com/stefanprodan/kubectl-kubesec/releases/download/0.1.0/kubectl-kubesec_0.1.0_`uname -s`_amd64.tar.gz | tar xzvf - -C ~/.kube/plugins/scan
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
