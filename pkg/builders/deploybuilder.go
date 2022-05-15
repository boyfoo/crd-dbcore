package builders

import (
	"bytes"
	"context"
	"fmt"
	dbconfigv1 "github.com/shenyisyn/dbcore/pkg/apis/dbconfig/v1"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"text/template"
)

type DeployBuilder struct {
	deploy *v1.Deployment
	client.Client
	config *dbconfigv1.DbConfig
}

//目前软件的 命名规则
func deployName(name string) string {
	return "dbcore-" + name
}

// NewDeployBuilder 根据模板内容构建出deploy
func NewDeployBuilder(config dbconfigv1.DbConfig, client client.Client) (*DeployBuilder, error) {
	dep := &v1.Deployment{}

	if err := client.Get(context.Background(), types.NamespacedName{
		Namespace: config.Namespace,
		Name:      deployName(config.Name),
	}, dep); err != nil { // 没取到数据
		fmt.Println("查找不到 " + config.Name)
		dep.Name = config.Name
		dep.Namespace = config.Namespace
		tpl, err := template.New("deploy").Parse(deptpl)
		if err != nil {
			return nil, err
		}
		fmt.Println("构建模板成功")
		var tplRet bytes.Buffer
		err = tpl.Execute(&tplRet, dep) // 构建模板内容
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal(tplRet.Bytes(), dep) // 模板内容映射到结构体
		if err != nil {
			return nil, err
		}
		fmt.Println("模板内容映射到结构体成功")
	}

	return &DeployBuilder{
		deploy: dep,
		Client: client,
		config: &config,
	}, nil
}

// apply 同步属性
func (d *DeployBuilder) apply() *DeployBuilder {
	*d.deploy.Spec.Replicas = int32(d.config.Spec.Replicas)
	return d
}

// 设置关系owner
func (d *DeployBuilder) setOwner() *DeployBuilder {
	d.deploy.OwnerReferences = append(d.deploy.OwnerReferences, metav1.OwnerReference{
		APIVersion: d.config.APIVersion,
		Kind:       d.config.Kind,
		Name:       d.config.Name,
		UID:        d.config.UID,
	})
	return d
}

func (d *DeployBuilder) Build(ctx context.Context) error {
	// 未创建时没有值 就创建一个
	if d.deploy.CreationTimestamp.IsZero() {
		d.apply().setOwner()
		fmt.Println("创建新的 " + d.deploy.Name)
		err := d.Create(ctx, d.deploy)
		if err != nil {
			return err
		}
	} else {
		// 存在就是更新
		patch := client.MergeFrom(d.deploy.DeepCopy())
		d.apply()
		err := d.Patch(ctx, d.deploy, patch)
		if err != nil {
			return err
		}
		//修改状态 deploy已准备的pod
		replicas := d.deploy.Status.ReadyReplicas
		d.config.Status.Replicas = replicas
		d.config.Status.Ready = fmt.Sprintf("%d/%d", replicas, d.config.Spec.Replicas)

		err = d.Client.Status().Update(ctx, d.config)
		if err != nil {
			return err
		}
	}
	return nil
}
