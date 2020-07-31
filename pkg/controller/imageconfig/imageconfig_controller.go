package imageconfig

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	githubv1 "github.com/shenkonghui/image-controller/pkg/apis/github/v1"
	app "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"sync"
)

var ImageConfigCache sync.Map

var log = logf.Log.WithName("controller_imageconfig")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new ImageConfig Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileImageConfig{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("imageconfig-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource ImageConfig
	err = c.Watch(&source.Kind{Type: &githubv1.ImageConfig{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	nr := r.(* ReconcileImageConfig)
	pre := &predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			return nr.DeploymentCreateEvent(e)
			// Ignore updates to CR status in which case metadata.Generation does not change
			//return e.MetaOld.GetGeneration() != e.MetaNew.GetGeneration()
		},
	}
	err = c.Watch(&source.Kind{Type: &app.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &githubv1.ImageConfig{},
	},pre)
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileImageConfig implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileImageConfig{}

// ReconcileImageConfig reconciles a ImageConfig object
type ReconcileImageConfig struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a ImageConfig object and makes changes based on the state read
// and what is in the ImageConfig.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileImageConfig) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling ImageConfig")

	instance := &githubv1.ImageConfig{}
	err := r.client.Get(context.TODO(),request.NamespacedName,instance)
	// 更新内存中的对象
	if err != nil && errors.IsNotFound(err) {
		log.V(0).Info(fmt.Sprintf("imageconfig[%s] not find, Delete from cache",request.Name))
		ImageConfigCache.Range(func(key, value interface{}) bool {
			ic, ok := value.(*githubv1.ImageConfig)
			if !ok{
				return false
			}

			if ic.Name == instance.Name &&
				ic.Namespace == instance.Namespace{
				ImageConfigCache.Delete(fmt.Sprintf("%s/%s",request.Namespace,ic.Spec.Repo))
			}
			return true
		})

	} else {
		//r.logwithRqa(instance, "update rqa", 0)
		ImageConfigCache.Store(fmt.Sprintf("%s/%s",request.Namespace,instance.Spec.Repo),instance)
	}
	// Pod already exists - don't requeue
	return reconcile.Result{}, nil
}
