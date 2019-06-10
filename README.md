# kubectl-kubesec

[![Build Status](https://travis-ci.org/controlplaneio/kubectl-kubesec.svg?branch=master)](https://travis-ci.org/controlplaneio/kubectl-kubesec)

This is a kubectl plugin for scanning Kubernetes pods, deployments, daemonsets and statefulsets with [kubesec.io](https://kubesec.io)

For the admission controller see [kubesec-webhook](https://github.com/controlplaneio/kubesec-webhook)

### Install

#### Install with krew

1. [Install krew](https://github.com/GoogleContainerTools/krew) plugin manager
   for kubectl.
2. Run `kubectl krew install kubesec-scan`.
3. Start using by running `kubectl kubesec-scan`.

#### Install with curl

For Kubernetes 1.12 or newer:

```bash
mkdir -p ~/.kube/plugins/scan && \
curl -sL https://github.com/controlplaneio/kubectl-kubesec/releases/download/1.0.0/kubectl-kubesec_1.0.0_`uname -s`_amd64.tar.gz | tar xzvf - -C ~/.kube/plugins/scan
mv ~/.kube/plugins/scan/scan ~/.kube/plugins/scan/kubectl-scan
export PATH=$PATH:~/.kube/plugins/scan
```

For Kubernetes older than 1.12:

```bash
mkdir -p ~/.kube/plugins/scan && \
curl -sL https://github.com/controlplaneio/kubectl-kubesec/releases/download/0.3.1/kubectl-kubesec_0.3.1_`uname -s`_amd64.tar.gz | tar xzvf - -C ~/.kube/plugins/scan
```

### Usage

Scan a Deployment:

```bash
kubectl scan -n kube-system deployment kubernetes-dashboard
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

Scan a DaemonSet:

```bash
kubectl scan -n weave daemonset weave-scope-agent
```

Result:

```
kubesec.io score: -51
-----------------
Critical
1. .spec .hostNetwork == true
Sharing the host's network namespace permits processes in the pod to communicate with processes bound to the host's loopback adapter
2. .spec .hostPID == true
Sharing the host's PID namespace allows visibility of processes on the host, potentially leaking information such as environment variables and configuration
3. containers[] .securityContext .privileged == true
Privileged containers can allow almost completely unrestricted host access
4. volumes[] .hostPath .path == /var/run/docker.sock
Mounting the docker.socket leaks information about other containers and can allow container breakout
-----------------
Advise
1. containers[] .securityContext .runAsUser -gt 10000
Run as a high-UID user to avoid conflicts with the host's user table
2. containers[] .securityContext .readOnlyRootFilesystem == true
An immutable root filesystem can prevent malicious binaries being added to PATH and increase attack cost
3. containers[] .securityContext .runAsNonRoot == true
Force the running image to run as a non-root user to ensure least privilege
4. containers[] .resources .limits .cpu
Enforcing CPU limits prevents DOS via resource exhaustion
5. containers[] .securityContext .capabilities .drop
Reducing kernel capabilities available to a container limits its attack surface
6. .metadata .annotations ."container.seccomp.security.alpha.kubernetes.io/pod"
Seccomp profiles set minimum privilege and secure against unknown threats
7. containers[] .securityContext .capabilities .drop | index("ALL")
Drop all capabilities and add only those required to reduce syscall attack surface
```

Scan a StatefulSet:

```bash
kubectl scan statefulset memcached
```

Result:

```
statefulset/memcached kubesec.io score 2
-----------------
Advise
1. .spec .volumeClaimTemplates[] .spec .accessModes | index("ReadWriteOnce")
2. containers[] .securityContext .runAsNonRoot == true
Force the running image to run as a non-root user to ensure least privilege
3. containers[] .securityContext .capabilities .drop
Reducing kernel capabilities available to a container limits its attack surface
4. containers[] .securityContext .readOnlyRootFilesystem == true
An immutable root filesystem can prevent malicious binaries being added to PATH and increase attack cost
5. containers[] .securityContext .runAsUser > 10000
Run as a high-UID user to avoid conflicts with the host's user table
```

Scan a Pod:

```bash
kubectl scan -n kube-system pod tiller-deploy-5c688d5f9b-ztjbt
```

Result:

```
pod/tiller-deploy-5c688d5f9b-ztjbt kubesec.io score 3
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
