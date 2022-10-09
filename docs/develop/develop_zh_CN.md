# 开发指南

## 本地测试

> 开发机为Mac，使用minikube

- 本地需装有[minikube](https://minikube.sigs.k8s.io/docs/start/)，驱动选择virtualbox。

### 构建&测试

- 使用minikube docker守护进程：eval $(minikube docker-env)

### 添加磁盘

- 默认启动的minikube只有一块磁盘（/dev/sda），添加磁盘步骤如下
  - 切换minikube目录：cd ~/.minikube/machines/minikube/minikube
  - 创建磁盘文件：VBoxManage createhd --filename mydisk00.vdi --size 40960
  - 查看device、port信息：cat minikube.vbox|grep HardDisk
  - 磁盘attach，此处device为上述信息中的device序号，port则在原基础上加1：VBoxManage storageattach minikube --storagectl "SATA" --port 2 --device 0 --type hdd --medium mydisk00.vdi
  - 登陆进minikube虚拟机：minikube ssh
  - 查看磁盘是否挂载成功：lsblk

### 更新lvm.proto文件

```bash
go get google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
export PATH="$PATH:$(go env GOPATH)/bin"
cd pkg/csi/lib
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    lvm.proto
```