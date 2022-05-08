### 初始化

`go get -u k8s.io/apimachinery@v0.20.5`

`go get -u k8s.io/client-go@v0.20.5`

添加对应的文件夹`pkg/apis/dbconfig/v1`，并新增文件`doc.go`，`register.go`，`types.go`，具体内容查询`v1-init`分支

### 生成代码

查看`go`路径：`go env | grep GOPATH`

将`v1-init`分支代码复制至`$GOPATH/src/github.com/shenyisyn/dbcore` 目录下，这个目录是`go.mod`内的`module`名称

下载生成器地址: https://github.com/kubernetes/code-generator/releases/tag/kubernetes-1.20.5

将下载好的生成器移动至`$GOPATH/src/k8s.io`下，并重命名为`code-generator`，并执行`go mod download`

在`$GOPATH/src/github.com/shenyisyn/dbcore`
目录下执行代码生成`$GOPATH/src/k8s.io/code-generator/generate-groups.sh all  github.com/shenyisyn/dbcore/pkg/client github.com/shenyisyn/dbcore/pkg/apis dbconfig:v1`

再将生成好的`v1/zz_generated.deepcopy.go`拷贝回来，将`client/`拷贝到`pkg`文件夹下

之所以这么麻烦，是因为生成器本身是为了生成`src`下的包而编写的，需要把自己的项目模拟成该形态

### 新增工具代码

新增 `pkg/k8sconfig`，内提供获取客户端的方法，目前代码截止查看`tag:v1.0.1`

### 发布自定义crd

发布crd: `kb apply -f crd/crd.yaml`

发布测试资源：`kb apply -f crd/crd-test.yaml`

查询资源: `kb get dc`

### 简单的使用client

参看`v1.0.2`的`main.go`文件

### 新增控制器管理基础版

在`go.mod`中新增`sigs.k8s.io/controller-runtime v0.9.6`

新增文件`controllers/DbConfigController.go`

新建控制器管理文件`k8sconfig/initManger.go`

`main.go`中调用初始化控制器管理

具体代码查看`tag:v1.0.3`

### 自定义字段验证

文档：https://kubernetes.io/zh/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/#validation

规范文档：https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.0.md

### 新增status子资源

见`tag:v1.0.4`文件`crd/crd.yaml`，修改`pkg/apis/dbconfig/v1/types.go`内`DbConfigStatus.Replicas`类型为`string`

### 新增字段打印和伸缩属性设置字段

查看`tag:v1.0.5`

在`crd.yaml`中新增`additionalPrinterColumns`参数添加打印字段

在`crd.yaml`中`subresources.scale`新增伸缩字段路径

