![Unit Tests](https://github.com/vultr/vultr-csi/workflows/Unit%20Tests/badge.svg)
# vultr-csi

The Container Storage Interface ([CSI](https://github.com/container-storage-interface/spec)) Driver for Vultr [Block Storage](https://www.vultr.com/products/block-storage/). This driver allows you to use Vultr Block Storage with your container orchestrator. We have tested this CSI on Kubernetes.

More information about the CSI and Kubernetes can be found: [CSI Spec](https://github.com/container-storage-interface/spec) and [Kubernetes CSI](https://kubernetes-csi.github.io/docs/example.html)


## Installation
### Requirements

- `--allow-privileged` must be enabled for the API server and kubelet

### Kubernetes secret

In order for the csi to work properly, you will need to deploy a [kubernetes secret](https://kubernetes.io/docs/concepts/configuration/secret/). To obtain a API key, please visit [API settings](https://my.vultr.com/settings/#settingsapi).  

The `secret.yml` definition is as follows. You can also find a copy of this yaml [here](docs/releases/secret.yml.tmp).
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: vultr-csi
  namespace: kube-system
stringData:
  # Replace the api-key with a proper value
  api-key: "VULTR_API_KEY"
```

To create this `secret.yml`, you must run the following

```sh
$ kubectl create -f secret.yml            
secret/vultr-csi created
```

### Deploying the CSI

To deploy the latest release of the CSI to your Kubernetes cluster, run the following:

`kubectl apply -f https://raw.githubusercontent.com/vultr/vultr-csi/master/docs/releases/latest.yml`

If you wish to deploy a specific version, you must replace `latest` with a proper release where `X.Y.Z` is the desired version:

`https://raw.githubusercontent.com/vultr/vultr-csi/master/docs/releases/vX.Y.Z.yml`

### Validating

The deployment will create a [Storage Class](https://kubernetes.io/docs/concepts/storage/storage-classes/) which will be used to create your volumes

```sh
$ kubectl get storageclass
NAME                            PROVISIONER              RECLAIMPOLICY   VOLUMEBINDINGMODE   ALLOWVOLUMEEXPANSION   AGE
vultr-block-storage (default)   vultrbs.csi.driver.com   Delete          Immediate           false                  131m
vultr-block-storage-retain      vultrbs.csi.driver.com   Retain          Immediate           false                  131m
```

To further validate the CSI, create a [PersistentVolumeClaim](https://kubernetes.io/docs/concepts/storage/persistent-volumes/)

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: csi-pvc
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
  storageClassName: vultr-block-storage
```

Now, take the yaml shown above and create a `pvc.yml` and run:

`kubectl create -f pvc.yml`

You can then check that you have a unattached volume on the Vultr dashboard. In addition, you can see that you have a `PersistentVolume` created by your Claim

```sh
$ kubectl get pv
NAME                   CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM             STORAGECLASS          REASON   AGE
pvc-2579a832202d4d07   10Gi       RWO            Delete           Bound    default/csi-pvc   vultr-block-storage            2s
``` 

Again, this volume is not attached to any node/pod yet. The volume will be attached to a node when a pod residing inside that node requests the specific volume.

Here is an example yaml of a pod request for the volume we just created.

```yaml
kind: Pod
apiVersion: v1
metadata:
  name: readme-app
spec:
  containers:
    - name: readme-app
      image: busybox
      volumeMounts:
      - mountPath: "/data"
        name: vultr-volume
      command: [ "sleep", "1000000" ]
  volumes:
    - name: vultr-volume
      persistentVolumeClaim:
        claimName: csi-pvc 
```
`kubectl create -f pod-volume.yml`

To get more information about the pod to ensure it is running and mounted, you can run the following

`kubectl describe po readme-app`

Now, let's add some data to the pod and validate that if we delete a pod and recreate a new pod which requests the same volume, the data still exists.

```sh
# Create a file
$ kubectl exec -it readme-app -- /bin/sh -c "touch /data/example"

# Delete the Pod
kubectl delete -f pod-volume.yml

# Recreate the pod with the same volume
kubectl create -f pod-volume.yml

# See that data on our volume still exists
$ kubectl exec -it readme-app -- /bin/sh -c "ls /data"   
```

## Examples
Some example yaml definitions can be found [here](docs/examples)

## Contributing Guidelines
If you are interested in improving or helping with vultr-csi, please feel free to open an issue or PR!
