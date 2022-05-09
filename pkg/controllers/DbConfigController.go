package controllers

import (
	"context"
	v1 "github.com/shenyisyn/dbcore/pkg/apis/dbconfig/v1"
	"github.com/shenyisyn/dbcore/pkg/builders"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type DbConfigController struct {
	client.Client
}

func (r *DbConfigController) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	config := &v1.DbConfig{}
	err := r.Get(ctx, req.NamespacedName, config)
	if err != nil {
		return reconcile.Result{}, err
	}
	builder, err := builders.NewDeployBuilder(*config, r.Client)
	if err != nil {
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

func NewDbConfigController() *DbConfigController {
	return &DbConfigController{}
}
