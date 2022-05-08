package main

import (
	"context"
	"fmt"
	v1 "github.com/shenyisyn/dbcore/pkg/client/clientset/versioned/typed/dbconfig/v1"
	"github.com/shenyisyn/dbcore/pkg/k8sconfig"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	client, err := v1.NewForConfig(k8sconfig.K8sRestConfig())
	if err != nil {
		fmt.Println(err)
		return
	}
	list, err := client.DbConfigs("default").List(context.Background(), metav1.ListOptions{})
	for _, item := range list.Items {
		fmt.Println(item.Name)
	}
}
