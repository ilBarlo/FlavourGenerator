package flavourgenerator

import (
	"context"
	"fmt"
	"log"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func StartController(cl client.Client) {
	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{})
	if err != nil {
		log.Fatalf("unable to create manager: %v", err)
	}

	kubeClient, err := kubernetes.NewForConfig(config.GetConfigOrDie())
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}

	//labelSelector := labels.Set{workerLabelKey: ""}

	// Create a new informer for pods
	podInformer := informers.NewSharedInformerFactoryWithOptions(
		kubeClient, time.Minute,
	).Core().V1().Pods().Informer()

	// Add event handler function to informer
	podInformer.AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: func(obj interface{}) bool {
			// Filter out all non-modify events
			if _, ok := obj.(metav1.Object); ok {
				return true
			}
			return false
		},
		Handler: cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				nodes, err := GetNodesResources(context.Background(), cl)
				if err != nil {
					log.Printf("Error getting nodes resources: %v", err)
					return
				}

				for _, node := range *nodes {
					flavours := splitResources(node)
					for _, flavour := range flavours {
						sendMessage(flavour, "flavours", natsURL)
						fmt.Printf("Flavour sent from node %s\n", flavour.Name)
					}
				}
			},
			DeleteFunc: func(obj interface{}) {
				nodes, err := GetNodesResources(context.Background(), cl)
				if err != nil {
					log.Printf("Error getting nodes resources: %v", err)
					return
				}

				for _, node := range *nodes {
					flavours := splitResources(node)
					for _, flavour := range flavours {
						sendMessage(flavour, "flavours", natsURL)
						fmt.Printf("Flavour sent from node %s\n", flavour.Name)
					}
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				nodes, err := GetNodesResources(context.Background(), cl)
				if err != nil {
					log.Printf("Error getting nodes resources: %v", err)
					return
				}

				for _, node := range *nodes {
					flavours := splitResources(node)
					for _, flavour := range flavours {
						sendMessage(flavour, "flavours", natsURL)
						fmt.Printf("Flavour sent from node %s\n", flavour.Name)
					}
				}
			},
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go podInformer.Run(ctx.Done())

	if err := mgr.Start(ctx); err != nil {
		log.Fatalf("unable to start manager: %v", err)
	}
}
