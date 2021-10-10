# kubectl-kubesec

[![Build Status](https://travis-ci.org/controlplaneio/kubectl-kubesec.svg?branch=master)](https://travis-ci.org/controlplaneio/kubectl-kubesec)

This is a kubectl plugin for scanning Kubernetes pods, deployments, daemonsets and statefulsets with [kubesec.io](https://kubesec.io). By default the plugin will send scan requests to the hosted version of [kubesec.io](https://kubesec.io). However, it is also possible to self host the scanning service and use that for scanning instead.

For the admission controller see [kubesec-webhook](https://github.com/controlplaneio/kubesec-webhook)

The latest release of this plugin is fully compatible with the API version V2 of kubesec documented at [kubesec.io](https://kubesec.io).

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

By default the plugin uses the hosted version of [kubesec.io](https://kubesec.io). However, you can run the hosted service locally. For example using docker:

```bash
## 
docker run -d -p 8080:8080 kubesec/kubesec:v2 http 8080
```

Scan a Deployment:

```bash
kubectl kubesec-scan -n kube-system deployment kubernetes-dashboard
# if you are running a self hosted version of kubese.io using docker then:
kubectl kubesec-scan -n kube-system deployment kubernetes-dashboard --url http://localhost:8080
```

Result:

```bash
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
kubectl kubesec-scan -n weave daemonset weave-scope-agent
# if you are running a self hosted version of kubese.io using then:
kubectl kubesec-scan -n weave daemonset weave-scope-agent --url http://localhost:8080
```

Result:

```bash
daemonset/weave-scope-agent kubesec.io score -54
-----------------
Critical
1. containers[] .securityContext .privileged == true
Privileged containers can allow almost completely unrestricted host access
2. .spec .hostNetwork
Sharing the host's network namespace permits processes in the pod to communicate with processes bound to the host's loopback adapter
3. .spec .hostPID
Sharing the host's PID namespace allows visibility of processes on the host, potentially leaking information such as environment variables and configuration
4. .spec .volumes[] .hostPath .path == "/var/run/docker.sock"
Mounting the docker.socket leaks information about other containers and can allow container breakout
```

Scan a StatefulSet:

```bash
kubectl kubesec-scan statefulset memcached
# if you are running a self hosted version of kubese.io then:
kubectl kubesec-scan statefulset memcached --url http://localhost:8080
```

Result:

```bash
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
kubectl kubesec-scan -n kube-system pod tiller-deploy-5c688d5f9b-ztjbt
# if you are running a self hosted version of kubese.io then:
kubectl kubesec-scan -n kube-system pod tiller-deploy-5c688d5f9b-ztjbt --url http://localhost:8080 
```

Result:

```bash
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
