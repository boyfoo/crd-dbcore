package k8sconfig

import (
	v1 "github.com/shenyisyn/dbcore/pkg/apis/dbconfig/v1"
	"github.com/shenyisyn/dbcore/pkg/controllers"
	"log"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	k8slog "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
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

	err = builder.ControllerManagedBy(mgr).For(&v1.DbConfig{}).Complete(controllers.NewDbConfigController())
	if err != nil {
		mgr.GetLogger().Error(err, "控制器新增失败")
		os.Exit(1)
	}

	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		mgr.GetLogger().Error(err, "Start error")
		os.Exit(1)
	}
}
