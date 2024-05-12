package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/willianpc/k8s-client-playground/handlers"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)

	config, err := clientcmd.BuildConfigFromFlags("", "/home/willian/rpi.kubeconfig")

	if err != nil {
		log.Fatal(err)
	}

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		log.Fatal(err)
	}

	// go monitor(clientset)
	go monitorInformer(clientset, quitChannel)

	<-quitChannel
}

func monitor(clientSet *kubernetes.Clientset) {
	ctx := context.Background()

	watcher, err := clientSet.AppsV1().Deployments("willian").Watch(ctx, metav1.ListOptions{})

	if err != nil {
		log.Fatal(err)
	}

	for ev := range watcher.ResultChan() {
		if d, ok := ev.Object.(*appsv1.Deployment); !ok {
			log.Fatal("supposed to be deployment", err)
		} else {
			log.Println("Deployment", d.Name, "Event type is", ev.Type)
		}
	}
}

func monitorInformer(clientSet *kubernetes.Clientset, quitChannel chan os.Signal) {
	ch := make(chan struct{}, 1)

	defer func() {
		<-quitChannel
		<-ch
	}()

	factory := informers.NewSharedInformerFactory(clientSet, time.Hour)

	di := factory.Apps().V1().Deployments()

	di.Informer().AddEventHandler(&handlers.DeploymentHandler{})

	factory.Start(ch)
}
