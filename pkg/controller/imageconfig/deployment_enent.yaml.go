package imageconfig

import (
	"context"
	"fmt"
	githubv1 "github.com/shenkonghui/image-controller/pkg/apis/github/v1"
	app "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"strings"
	"time"
)

func (r *ReconcileImageConfig) DeploymentCreateEvent (e event.UpdateEvent) bool{
	instance, ok := e.MetaNew.(* app.Deployment)
	if !ok {
		return false
	}
	time.Sleep(time.Duration(3)*time.Second)

	err := r.client.Get(context.TODO(), types.NamespacedName{
		Namespace: instance.Namespace,
		Name: instance.Name}, instance)
	if err != nil{
		log.V(-1).Info(err.Error())
		return true
	}

	containers := instance.Spec.Template.Spec.Containers
	index := 0
	for _, c := range containers{
		s := getImageAdd(c.Image)
		o, ok := ImageConfigCache.Load(fmt.Sprintf("%s/%s",instance.Namespace,s[0]))
		if !ok{
			index ++
			continue
		}
		ic := o.(* githubv1.ImageConfig)
		if s[0] == ic.Spec.Repo{
			c.Image = fmt.Sprintf("%s/%s/%s",ic.Spec.NewRepo,s[1],s[2])
		}
		instance.Spec.Template.Spec.Containers[index] = c
		index ++
	}


	err = r.client.Update(context.TODO(), instance,&client.UpdateOptions{})
	if err != nil{
		log.V(-1).Info(err.Error())
		return true
	}
	return false

}

func getImageAdd(image string) []string{
	res := strings.Split(image,"/")
	if len(res) != 3 {
		resnew := make([]string, 3)
		resnew = append(resnew, "docker.io")
		resnew = append(resnew, "library")
		resnew = append(res, res[0])
		return resnew
	}
	return res
}