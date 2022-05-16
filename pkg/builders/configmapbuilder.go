package builders

import (
	"bytes"
	"context"
	"fmt"
	configv1 "github.com/shenyisyn/dbcore/pkg/apis/dbconfig/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"log"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"text/template"
)

type ConfigMapBuilder struct {
	cm     *corev1.ConfigMap
	config *configv1.DbConfig
	client.Client
}

func NewConfigMapBuilder(config *configv1.DbConfig, client client.Client) (*ConfigMapBuilder, error) {
	cm := &corev1.ConfigMap{}
	err := client.Get(context.Background(), types.NamespacedName{
		Namespace: config.Namespace, Name: deployName(config.Name), //这里做了改动
	}, cm)
	if err != nil { //没取到
		cm.Name, cm.Namespace = deployName(config.Name), config.Namespace
		cm.Data = make(map[string]string)
	}
	return &ConfigMapBuilder{cm: cm, Client: client, config: config}, nil
}

//同步属性
func (this *ConfigMapBuilder) setOwner() *ConfigMapBuilder {
	this.cm.OwnerReferences = append(this.cm.OwnerReferences,
		metav1.OwnerReference{
			APIVersion: this.config.APIVersion,
			Kind:       this.config.Kind,
			Name:       this.config.Name,
			UID:        this.config.UID,
		})
	return this
}
func (this *ConfigMapBuilder) apply() *ConfigMapBuilder {
	// Delims 设置分隔符，这边是测试一下用其他类型的分隔符
	tpl, err := template.New("app.yml").Delims("[[", "]]").Parse(cmtpl)
	if err != nil {
		log.Println("configmap apply error: " + err.Error())
		return this
	}
	var tplRet bytes.Buffer
	err = tpl.Execute(&tplRet, this.config.Spec)
	if err != nil {
		log.Println(err)
		return this
	}
	this.cm.Data["app.yml"] = tplRet.String()
	fmt.Println("this.cm.Data")
	fmt.Println(this.cm.Data)
	return this
}

func (this *ConfigMapBuilder) Build(ctx context.Context) error {
	if this.cm.CreationTimestamp.IsZero() {
		this.apply().setOwner() //同步  所需要的属性 如 副本数 , 并且设置OwnerReferences
		err := this.Create(ctx, this.cm)
		if err != nil {
			return err
		}
	} else {
		patch := client.MergeFrom(this.cm.DeepCopy())
		this.apply() //同步  所需要的属性 如 副本数
		err := this.Patch(ctx, this.cm, patch)
		if err != nil {
			return err
		}
	}
	//更新 是课后作业
	return nil
}
