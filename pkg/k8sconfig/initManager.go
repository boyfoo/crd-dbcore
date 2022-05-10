package k8sconfig

import (
	v1 "github.com/shenyisyn/dbcore/pkg/apis/dbconfig/v1"
	"github.com/shenyisyn/dbcore/pkg/controllers"
	appsv1 "k8s.io/api/apps/v1"
	"log"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	k8slog "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

func InitManager() {
	k8slog.SetLogger(zap.New())

	mgr, err := manager.New(K8sRestConfig(), manager.Options{
		Logger: k8slog.Log.WithName("dbcore"),
	})
	if err != nil {
		log.Fatal("创建管理器失败", err.Error())
	}

	err = v1.AddToScheme(mgr.GetScheme())
	if err != nil {
		mgr.GetLogger().Error(err, "新增 scheme失败")
		os.Exit(1)
	}

	controller := controllers.NewDbConfigController()
	err = builder.ControllerManagedBy(mgr).For(&v1.DbConfig{}).
		Watches(&source.Kind{Type: &appsv1.Deployment{}}, handler.Funcs{
			DeleteFunc: controller.OnDelete,
		}).
		Complete(controller)
	if err != nil {
		mgr.GetLogger().Error(err, "控制器新增失败")
		os.Exit(1)
	}

	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		mgr.GetLogger().Error(err, "Start error")
		os.Exit(1)
	}
}
