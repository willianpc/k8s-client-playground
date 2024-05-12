package handlers

import (
	"log"

	appsv1 "k8s.io/api/apps/v1"
)

const (
	HandlerAdd = iota
	HandlerUpdate
	HandlerDelete
)

type DeploymentHandler struct{}

func reportDeployment(obj interface{}, handlerType string) {
	d, ok := obj.(*appsv1.Deployment)

	if !ok {
		log.Println("object is not deployment", obj)
		return
	}

	log.Println("Deployment", d.Name, "with operation", handlerType)
}

func (dh *DeploymentHandler) OnAdd(obj interface{}, isInInitialList bool) {
	reportDeployment(obj, "add")
}

func (dh *DeploymentHandler) OnUpdate(oldObj, newObj interface{}) {
	reportDeployment(oldObj, "update")
}

func (dh *DeploymentHandler) OnDelete(obj interface{}) {
	reportDeployment(obj, "delete")
}
