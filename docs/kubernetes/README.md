## Block Types

As of version `v0.6.0` there are two available Block Storage types that you can deploy against your Kubernetes cluster.

### NVME Block Type

This Block type uses the storage class `vultr-block-storage`. It has a minimum deployment size of 1gb and a maximum of 100tb.

It is currently available in the following regions:

- Atlanta
- Amsterdam
- Bangalore
- Chicago
- Los Angeles
- London
- New Jersey
- Seattle
- Singapore
- Sydney
- Tokyo

### HDD Block Type

This Block type uses the storage class `vultr-block-storage-hdd`. It has a minimum deployment size of 40gb and a maximum of 40tb.

It is currently available in every region **except** for the following:

- Sydney

## Installation

### Requirements

- `--allow-privileged` must be enabled for the API server and kubelet
- Kubernetes 1.25 or newer is required for the bundled external-snapshotter

### Kubernetes secret

In order for the csi to work properly, you will need to deploy a
[kubernetes secret](https://kubernetes.io/docs/concepts/configuration/secret/).
To obtain a API key, please visit
[API settings](https://my.vultr.com/settings/#settingsapi).

The `secret.yml` definition is as follows. You can also find a copy of this yaml
[here](../releases/secret.yml.tmp).

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

Snapshots require the Kubernetes snapshot CRDs and the cluster-wide snapshot
controller. Some managed Kubernetes distributions install these components for
you. Check before installing them:

```sh
kubectl get crd volumesnapshots.snapshot.storage.k8s.io \
  volumesnapshotcontents.snapshot.storage.k8s.io \
  volumesnapshotclasses.snapshot.storage.k8s.io
kubectl -n kube-system get deployment snapshot-controller
```

If they are absent, install the v8.4.0 snapshot CRDs and the controller bundled
with this repository:

```sh
kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/v8.4.0/client/config/crd/snapshot.storage.k8s.io_volumesnapshotclasses.yaml
kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/v8.4.0/client/config/crd/snapshot.storage.k8s.io_volumesnapshotcontents.yaml
kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/v8.4.0/client/config/crd/snapshot.storage.k8s.io_volumesnapshots.yaml
kubectl apply -f https://raw.githubusercontent.com/vultr/vultr-csi/master/docs/kubernetes/snapshot-controller.yml
```

To deploy the latest release of the CSI to your Kubernetes cluster, run the
following:

`kubectl apply -f https://raw.githubusercontent.com/vultr/vultr-csi/master/docs/releases/latest.yml`

If you wish to deploy a specific version, you must replace `latest` with a
proper release where `X.Y.Z` is the desired version:

`https://raw.githubusercontent.com/vultr/vultr-csi/master/docs/releases/vX.Y.Z.yml`

### Validating

The deployment will create a
[Storage Class](https://kubernetes.io/docs/concepts/storage/storage-classes/)
which will be used to create your volumes

```sh
$ kubectl get storageclass
NAME                         PROVISIONER           RECLAIMPOLICY   VOLUMEBINDINGMODE   ALLOWVOLUMEEXPANSION
vultr-block-storage          block.csi.vultr.com   Delete          Immediate           true
vultr-block-storage-retain   block.csi.vultr.com   Retain          Immediate           true
```

To further validate the CSI, create a
[PersistentVolumeClaim](https://kubernetes.io/docs/concepts/storage/persistent-volumes/)

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

You can then check that you have a unattached volume on the Vultr dashboard. In
addition, you can see that you have a `PersistentVolume` created by your Claim

```sh
$ kubectl get pv
NAME                   CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM             STORAGECLASS          REASON   AGE
pvc-2579a832202d4d07   10Gi       RWO            Delete           Bound    default/csi-pvc   vultr-block-storage            2s
```

Again, this volume is not attached to any node/pod yet. The volume will be
attached to a node when a pod residing inside that node requests the specific
volume.

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

To get more information about the pod to ensure it is running and mounted, you
can run the following

`kubectl describe po readme-app`

Now, let's add some data to the pod and validate that if we delete a pod and
recreate a new pod which requests the same volume, the data still exists.

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

## Snapshots

Snapshots and restores are supported for Vultr Block Storage volumes. VFS
volumes cannot be snapshotted. The release manifest creates `Delete` and
`Retain` snapshot classes named `vultr-block-storage` and
`vultr-block-storage-retain`.

Create a snapshot of the example `csi-pvc` claim and wait until it is ready:

```sh
kubectl apply -f https://raw.githubusercontent.com/vultr/vultr-csi/master/docs/kubernetes/examples/snapshot.yml
kubectl wait --for=jsonpath='{.status.readyToUse}'=true \
  volumesnapshot/csi-snapshot --timeout=10m
```

Restore that snapshot into a new volume:

```sh
kubectl apply -f https://raw.githubusercontent.com/vultr/vultr-csi/master/docs/kubernetes/examples/restore-pvc.yml
kubectl get pvc csi-restored-pvc
```

The restored PVC must request at least as much capacity as the source snapshot.
Deleting a snapshot that uses `vultr-block-storage` also deletes the Vultr
snapshot. The `vultr-block-storage-retain` class preserves it.

## Examples

Some example yaml definitions can be found [here](examples)
