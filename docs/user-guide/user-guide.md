# User guide

This is the user guide for Open-Local.

## Requirements

- Linux distributions based on Red Hat and Debian
- Kubernetes v1.20+
- Helm v3.0+
- [lvm2](https://en.wikipedia.org/wiki/Logical_Volume_Manager_(Linux))
- At least one block device

## Configuring

To create a Volume Group by Open Local automatically, edit {{ .Values.agent.device }} in helm/values.yaml to change the name of the block device before deploying Open-Local in Kubernetes. It is recommended that **each node** provide a block device for Open-Local.

The name of the block device can be /dev/sdb, /dev/vdc, /dev/sdd2, etc.

## Deploying Open-Local

Install Open-Local with Helm 3, use these commands:

```bash
# helm install open-local ./helm
```

Confirm that the deployment succeeded, use these commands:

```bash
# kubectl get po -nkube-system  -l app=open-local
```

The following output should be displayed:

```bash
NAME                                             READY   STATUS      RESTARTS   AGE
open-local-agent-6zmkb                           3/3     Running     0          28s
open-local-csi-provisioner-6dbb7c459c-mcp9l      1/1     Running     0          28s
open-local-csi-resizer-57cfd85df7-x44zg          1/1     Running     0          28s
open-local-csi-snapshotter-689b6bbcfc-wwc57      1/1     Running     0          28s
open-local-init-job-2wvbs                        0/1     Completed   0          28s
open-local-init-job-bw8nh                        0/1     Completed   0          28s
open-local-init-job-frdxp                        0/1     Completed   0          28s
open-local-scheduler-extender-7d5cf688b6-pr426   1/1     Running     0          28s
open-local-snapshot-controller-d6f78754-czfnw    1/1     Running     0          28s
```

## Local storage pool management

Open-Local will create custom resources(nodelocalstorage) to report the storage information of each node in the cluster.

```bash
# kubectl get nodelocalstorage
NAME       STATE       PHASE     AGENTUPDATEAT   SCHEDULERUPDATEAT   SCHEDULERUPDATESTATUS
minikube   DiskReady   Running   30s             0s
```

Use following command to check if volume group is created and managed successfully by Open-Local:

```bash
# kubectl get nodelocalstorage -ojson minikube|jq .status.filteredStorageInfo
{
  "updateStatusInfo": {
    "lastUpdateTime": "2021-09-23T15:37:21Z",
    "updateStatus": "accepted"
  },
  "volumeGroups": [
    "open-local-pool-0"
  ]
}
```

## Dynamic volume provisioning

Open-Local has storageclasses as following:

```bash
NAME                    PROVISIONER                RECLAIMPOLICY   VOLUMEBINDINGMODE      ALLOWVOLUMEEXPANSION   AGE
open-local-device-hdd   local.csi.aliyun.com        Delete          WaitForFirstConsumer   false                  6h56m
open-local-device-ssd   local.csi.aliyun.com        Delete          WaitForFirstConsumer   false                  6h56m
open-local-lvm          local.csi.aliyun.com        Delete          WaitForFirstConsumer   true                   6h56m
```

Create a Pod that uses Open-Local volumes by running this command:

```bash
# kubectl apply -f ./example/lvm/sts-nginx.yaml
```

Check status of Pod/PVC/PV

```bash
# kubectl get pod
NAME          READY   STATUS    RESTARTS   AGE
nginx-lvm-0   1/1     Running   0          3m5s
# kubectl get pvc
NAME               STATUS   VOLUME                                       CAPACITY   ACCESS MODES   STORAGECLASS     AGE
html-nginx-lvm-0   Bound    local-52f1bab4-d39b-4cde-abad-6c5963b47761   5Gi        RWO            open-local-lvm   104s
# kubectl get pv
NAME                                         CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM                      STORAGECLASS    AGE
local-52f1bab4-d39b-4cde-abad-6c5963b47761   5Gi        RWO            Delete           Bound    default/html-nginx-lvm-0   open-local-lvm  2m4s
kubectl describe pvc html-nginx-lvm-0
Name:          html-nginx-lvm-0
Namespace:     default
StorageClass:  open-local-lvm
Status:        Bound
Volume:        local-52f1bab4-d39b-4cde-abad-6c5963b47761
Labels:        app=nginx-lvm
Annotations:   pv.kubernetes.io/bind-completed: yes
               pv.kubernetes.io/bound-by-controller: yes
               volume.beta.kubernetes.io/storage-provisioner: local.csi.aliyun.com
               volume.kubernetes.io/selected-node: minikube
Finalizers:    [kubernetes.io/pvc-protection]
Capacity:      5Gi
Access Modes:  RWO
VolumeMode:    Filesystem
Mounted By:    nginx-lvm-0
Events:
  Type    Reason                 Age                From                                                               Message
  ----    ------                 ----               ----                                                               -------
  Normal  WaitForFirstConsumer   11m                persistentvolume-controller                                        waiting for first consumer to be created before binding
  Normal  ExternalProvisioning   11m (x2 over 11m)  persistentvolume-controller                                        waiting for a volume to be created, either by external provisioner "local.csi.aliyun.com" or manually created by system administrator
  Normal  Provisioning           11m (x2 over 11m)  local.csi.aliyun.com_minikube_c4e4e0b8-4bac-41f7-88e4-149dba5bc058  External provisioner is provisioning volume for claim "default/html-nginx-lvm-0"
  Normal  ProvisioningSucceeded  11m (x2 over 11m)  local.csi.aliyun.com_minikube_c4e4e0b8-4bac-41f7-88e4-149dba5bc058  Successfully provisioned volume local-52f1bab4-d39b-4cde-abad-6c5963b47761
```

## Volume expansion

Modify the requested spec.resources.requests.storage of the PVC

```bash
# kubectl patch pvc html-nginx-lvm-0 -p '{"spec":{"resources":{"requests":{"storage":"20Gi"}}}}'
```

Check status of PVC/PV

```bash
# kubectl get pvc
NAME                    STATUS   VOLUME                                       CAPACITY   ACCESS MODES   STORAGECLASS     AGE
html-nginx-lvm-0        Bound    local-52f1bab4-d39b-4cde-abad-6c5963b47761   20Gi       RWO            open-local-lvm   7h4m
# kubectl get pv
NAME                                         CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM                           STORAGECLASS     REASON   AGE
local-52f1bab4-d39b-4cde-abad-6c5963b47761   20Gi       RWO            Delete           Bound    default/html-nginx-lvm-0        open-local-lvm            7h4m
```

## Volume snapshot

Open-Local has volumesnapshotclass as following:

```bash
NAME             DRIVER                DELETIONPOLICY   AGE
open-local-lvm   local.csi.aliyun.com   Delete           20m
```

Create a VolumeSnapshot

```bash
# kubectl apply -f example/lvm/snapshot.yaml
volumesnapshot.snapshot.storage.k8s.io/new-snapshot-test created
# kubectl get volumesnapshot
NAME                READYTOUSE   SOURCEPVC          SOURCESNAPSHOTCONTENT   RESTORESIZE   SNAPSHOTCLASS    SNAPSHOTCONTENT                                    CREATIONTIME   AGE
new-snapshot-test   true         html-nginx-lvm-0                           1863          open-local-lvm   snapcontent-815def28-8979-408e-86de-1e408033de65   19s            19s
# kubectl get volumesnapshotcontent
NAME                                               READYTOUSE   RESTORESIZE   DELETIONPOLICY   DRIVER                VOLUMESNAPSHOTCLASS   VOLUMESNAPSHOT      AGE
snapcontent-815def28-8979-408e-86de-1e408033de65   true         1863          Delete           local.csi.aliyun.com   open-local-lvm        new-snapshot-test   48s
```

Create a Pod that uses volume pre-populated with data from snapshots:

```bash
# kubectl apply -f example/lvm/sts-nginx-snap.yaml
service/nginx-lvm-snap created
statefulset.apps/nginx-lvm-snap created
# kubectl get po -l app=nginx-lvm-snap
NAME               READY   STATUS    RESTARTS   AGE
nginx-lvm-snap-0   1/1     Running   0          46s
# kubectl get pvc -l app=nginx-lvm-snap
NAME                    STATUS   VOLUME                                       CAPACITY   ACCESS MODES   STORAGECLASS     AGE
html-nginx-lvm-snap-0   Bound    local-1c69455d-c50b-422d-a5c0-2eb5c7d0d21b   4Gi        RWO            open-local-lvm   2m11s
# kubectl describe pvc html-nginx-lvm-snap-0
Name:          html-nginx-lvm-snap-0
Namespace:     default
StorageClass:  open-local-lvm
Status:        Bound
Volume:        local-1c69455d-c50b-422d-a5c0-2eb5c7d0d21b
Labels:        app=nginx-lvm-snap
Annotations:   pv.kubernetes.io/bind-completed: yes
               pv.kubernetes.io/bound-by-controller: yes
               volume.beta.kubernetes.io/storage-provisioner: local.csi.aliyun.com
               volume.kubernetes.io/selected-node: minikube
Finalizers:    [kubernetes.io/pvc-protection]
Capacity:      4Gi
Access Modes:  RWO
VolumeMode:    Filesystem
DataSource:
  APIGroup:  snapshot.storage.k8s.io
  Kind:      VolumeSnapshot
  Name:      new-snapshot-test
Mounted By:  nginx-lvm-snap-0
Events:
  Type    Reason                 Age                    From                                                               Message
  ----    ------                 ----                   ----                                                               -------
  Normal  WaitForFirstConsumer   2m37s                  persistentvolume-controller                                        waiting for first consumer to be created before binding
  Normal  ExternalProvisioning   2m37s                  persistentvolume-controller                                        waiting for a volume to be created, either by external provisioner "local.csi.aliyun.com" or manually created by system administrator
  Normal  Provisioning           2m37s (x2 over 2m37s)  local.csi.aliyun.com_minikube_c4e4e0b8-4bac-41f7-88e4-149dba5bc058  External provisioner is provisioning volume for claim "default/html-nginx-lvm-snap-0"
  Normal  ProvisioningSucceeded  2m37s (x2 over 2m37s)  local.csi.aliyun.com_minikube_c4e4e0b8-4bac-41f7-88e4-149dba5bc058  Successfully provisioned volume local-1c69455d-c50b-422d-a5c0-2eb5c7d0d21b
```

## Raw block volume

Open-Local also supports that the created storage volume will appear in the container as a block device (in this example, the block device is in the container /dev/sdd path):

```bash
# kubectl apply -f ./example/lvm/sts-block.yaml
```

Check status of Pod/PVC/PV:

```bash
# kubectl get pod
NAME                READY   STATUS    RESTARTS   AGE
nginx-lvm-block-0   1/1     Running   0          25s
# kubectl get pvc
NAME                     STATUS   VOLUME                                       CAPACITY   ACCESS MODES   STORAGECLASS     AGE
html-nginx-lvm-block-0   Bound    local-b048c19a-fe0b-455d-9f25-b23fdef03d8c   5Gi        RWO            open-local-lvm   36s
# kubectl get pv
NAME                                         CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM                            STORAGECLASS     REASON   AGE
local-b048c19a-fe0b-455d-9f25-b23fdef03d8c   5Gi        RWO            Delete           Bound    default/html-nginx-lvm-block-0   open-local-lvm            53s
# kubectl describe pvc html-nginx-lvm-block-0
Name:          html-nginx-lvm-block-0
Namespace:     default
StorageClass:  open-local-lvm
Status:        Bound
Volume:        local-b048c19a-fe0b-455d-9f25-b23fdef03d8c
Labels:        app=nginx-lvm-block
Annotations:   pv.kubernetes.io/bind-completed: yes
               pv.kubernetes.io/bound-by-controller: yes
               volume.beta.kubernetes.io/storage-provisioner: local.csi.aliyun.com
               volume.kubernetes.io/selected-node: izrj96fgmgzcvhtz2vkrgez
Finalizers:    [kubernetes.io/pvc-protection]
Capacity:      5Gi
Access Modes:  RWO
VolumeMode:    Block
Mounted By:    nginx-lvm-block-0
Events:
  Type    Reason                 Age                From                                                                               Message
  ----    ------                 ----               ----                                                                               -------
  Normal  WaitForFirstConsumer   72s                persistentvolume-controller                                                        waiting for first consumer to be created before binding
  Normal  Provisioning           72s                local.csi.aliyun.com_iZrj96fgmgzcvhtz2vkrgeZ_f2b69212-7103-4f9a-a6c4-179f37036ef0  External provisioner is provisioning volume for claim "default/html-nginx-lvm-block-0"
  Normal  ExternalProvisioning   72s (x2 over 72s)  persistentvolume-controller                                                        waiting for a volume to be created, either by external provisioner "local.csi.aliyun.com" or manually created by system administrator
  Normal  ProvisioningSucceeded  72s                local.csi.aliyun.com_iZrj96fgmgzcvhtz2vkrgeZ_f2b69212-7103-4f9a-a6c4-179f37036ef0  Successfully provisioned volume local-b048c19a-fe0b-455d-9f25-b23fdef03d8c
```