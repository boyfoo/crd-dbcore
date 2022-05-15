package controllers

import (
	"context"
	"fmt"
	v1 "github.com/shenyisyn/dbcore/pkg/apis/dbconfig/v1"
	"github.com/shenyisyn/dbcore/pkg/builders"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	Kind            = "DbConfig"
	GroupApiVersion = "api.jtthink.com/v1"
)

type DbConfigController struct {
	client.Client
}

func (r *DbConfigController) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	config := &v1.DbConfig{}
	err := r.Get(ctx, req.NamespacedName, config)
	if err != nil {
		fmt.Println("err := r.Get(ctx, req.NamespacedName, config)")
		return reconcile.Result{}, err
	}
	builder, err := builders.NewDeployBuilder(*config, r.Client)
	if err != nil {
		fmt.Println("builder, err := builders.NewDeployBuilder(*config, r.Client)")
		return reconcile.Result{}, err
	}
	err = builder.Build(ctx)
	if err != nil {
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

func (r *DbConfigController) InjectClient(c client.Client) error {
	r.Client = c
	return nil
}

func (r *DbConfigController) OnDelete(event event.DeleteEvent, limitingInterface workqueue.RateLimitingInterface) {
	// 监听的是deploy资源，当被删除的资源是所属对应的自定义资源Dbconfig时，触发 Reconcile 方法
	for _, ref := range event.Object.GetOwnerReferences() {
		if ref.Kind == Kind && ref.APIVersion == GroupApiVersion {
			limitingInterface.Add(reconcile.Request{
				NamespacedName: types.NamespacedName{Name: ref.Name, Namespace: event.Object.GetNamespace()},
			})
		}
	}
}

func (r *DbConfigController) OnUpdate(event event.UpdateEvent, limitingInterface workqueue.RateLimitingInterface) {
	for _, ref := range event.ObjectNew.GetOwnerReferences() {
		if ref.Kind == Kind && ref.APIVersion == GroupApiVersion {
			limitingInterface.Add(reconcile.Request{
				NamespacedName: types.NamespacedName{Name: ref.Name, Namespace: event.ObjectNew.GetNamespace()},
			})
		}
	}
}

func NewDbConfigController() *DbConfigController {
	return &DbConfigController{}
}
