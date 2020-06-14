package imageconfig

import (
	app "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

func (r *ReconcileImageConfig) DeploymentCreateEvent (e event.UpdateEvent) bool{
	instance, ok := e.MetaNew.(* app.Deployment)
	if !ok {
		return false
	}
	containers := instance.Spec.Template.Spec.Containers
	for _, c := range containers{

	}

}