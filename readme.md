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

测试结束，删除测试资源`kb delete -f crd/crd-test.yaml`


