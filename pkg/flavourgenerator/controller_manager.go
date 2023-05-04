package flavourgenerator

import (
	"context"
	"fmt"
	"os"

	corev1 "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// StartController by creating a new Mangaer, Reconciler and Controller instance
func StartController(cl client.Client) {
	// Create a new Manager instance
	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Create a new Reconciler instance
	r := &PodReconciler{
		Client: cl,
		Scheme: scheme,
	}

	// Create a new Controller instance
	c, err := ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Pod{}).
		Build(r)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Start the Manager
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Start the Controller
	if err := c.Start(ctrl.SetupSignalHandler()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Reconcile reads that state of the cluster for a Pod object and makes changes based on the state read
// and what is in the Pod.Spec
func (r *PodReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	nodes, err := GetNodesResources(ctx, r.Client)
	utilruntime.Must(err)

	for _, node := range *nodes {
		sendMessage(node, "metrics", "amqp://guest:guest@localhost:5672/")
		fmt.Printf("Metrics sent from node %s\n", node.Name)
	}
	return ctrl.Result{}, nil
}
