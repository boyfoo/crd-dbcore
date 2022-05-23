### 初始化

`go get -u k8s.io/apimachinery@v0.20.5`

`go get -u k8s.io/client-go@v0.20.5`

添加对应的文件夹`pkg/apis/dbconfig/v1`，并新增文件`doc.go`，`register.go`，`types.go`，具体内容查询`v1-init`分支

### 生成基础代码

#### 方式1

查看`go`路径：`go env | grep GOPATH`

将`v1-init`分支代码复制至`$GOPATH/src/github.com/shenyisyn/dbcore` 目录下，这个目录是`go.mod`内的`module`名称

下载生成器地址: https://github.com/kubernetes/code-generator/releases/tag/kubernetes-1.20.5

将下载好的生成器移动至`$GOPATH/src/k8s.io`下，并重命名为`code-generator`，并执行`go mod download`

在`$GOPATH/src/github.com/shenyisyn/dbcore`
目录下执行代码生成`$GOPATH/src/k8s.io/code-generator/generate-groups.sh all  github.com/shenyisyn/dbcore/pkg/client github.com/shenyisyn/dbcore/pkg/apis dbconfig:v1`

再将生成好的`v1/zz_generated.deepcopy.go`拷贝回来，将`client/`拷贝到`pkg`文件夹下

之所以这么麻烦，是因为生成器本身是为了生成`src`下的包而编写的，需要把自己的项目模拟成该形态

#### 方式2

下载生成器源码 `git clone git@github.com:kubernetes/code-generator.git`，切换到对于版本的分支`git checkout v0.20.2`

执行安装工具脚本 `go install ./cmd/{client-gen,deepcopy-gen,informer-gen,lister-gen}`

### 新增工具代码

新增 `pkg/k8sconfig`，内提供获取客户端的方法，目前代码截止查看`tag:v1.0.1`

### 发布自定义crd至k8s

发布crd: `kb apply -f crd/crd.yaml`

发布测试资源：`kb apply -f crd/crd-test.yaml`

查询资源: `kb get dc`

### 简单的使用client实例

参看`v1.0.2`的`main.go`文件

### 新增控制器管理基础版

在`go.mod`中新增`sigs.k8s.io/controller-runtime v0.9.6`

新增文件`controllers/DbConfigController.go`

新建控制器管理文件`k8sconfig/initManger.go`

`main.go`中调用初始化控制器管理

具体代码查看`tag:v1.0.3`

### 自定义资源字段验证

文档：https://kubernetes.io/zh/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/#validation

规范文档：https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.0.md

### 新增status子资源(status: 资源实际状态)

见`tag:v1.0.4`文件`crd/crd.yaml`，修改`pkg/apis/dbconfig/v1/types.go`内`DbConfigStatus.Replicas`类型为`string`

### 新增字段打印内容和伸缩属性设置

查看`tag:v1.0.5`

在`crd.yaml`中新增`additionalPrinterColumns`参数添加打印字段

在`crd.yaml`中`subresources.scale`新增伸缩字段路径

### 控制器根据自定义资源创建附属deploy

创建模板构建器文件`builders/deploybuilder.go`，修改`controllers/DbConfigController.go`文件代码，调用构建器根据模板`builders/deptpl.go`部署`deploy`

代码见`tag:v1.0.6`

### 控制器新增根据yaml修改功能

修改文件`builders/deploybuilder.go`内 `Build`方法和`apply`方法

代码见`tag:v1.0.7`

### 删除自定义资源触发删除附属子资源

当删除`DBconfig`资源时自动删除创建的`deploy`子资源，需要在创建`deploy`时设置`OwnerReferences`，就可以自动删除子资源

见`tag:v1.0.8`的`deploybuilder.go`文件下的`setOwner`方法

### 删除附属资源重新创建

在`ininManager`新增监听对应资源的删除`deploy`触发事件，对应函数判断该被删除的`deploy`是否是`DbConfig`的附属子资源

如果是重新触发`Reconcile`事件，目前是让资源从新创建`deploy`

代码见`tag:v1.0.9`

### 更新状态触发和设置命令行打印字段

修改`crd.yaml`文件内的`schema.status`子资源，新增`ready`字段和修改`replicas`类型

修改`crd.yaml`文件内`additionalPrinterColumns`下的`Ready`展示字段

修改`v1/types.go`内`DbConfigStatus`定义的格式，新增和修改属性与`crd.yaml`相同

更新`crd`:`k apply -f crd/crd.yaml`

修改`deptpl.go`新增一个休眠测试的`init`容器

修改`initManager.go`新增`UpdateFunc`监听函数，监听状态变化

见`tag:v1.0.10`

### 设置configmap

创建`configmap` 映射进`pod`中

并且运行`k apply -f mysql.yaml`，创建`mysql`，并创建库`test`

见 `tag:v1.1.0`

### 修改configmap内容重新隐射

将`configmap`计算成`md5`，每次更新更新`md5`，当值不一样时`deploy`重自动滚动更新，见`ConfigMapBuilder.DataKey`字段的使用

