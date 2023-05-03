package main

import (
	"context"
	"flag"
	"fmt"

	flavGenerator "github.com/ilbarlo/flavourGenerator/pkg/flavourgenerator"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/klog/v2"
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()
	ctx := context.Background()

	cl, err := flavGenerator.GetKClient(ctx)
	utilruntime.Must(err)

	str, err := flavGenerator.GetNodesResources(ctx, cl)
	utilruntime.Must(err)

	for _, node := range *str {
		fmt.Println(node.ResourceMetrics)
	}

	// klog.Fatalf("unable to start the server: %s", manager.SetupRouterAndServeHTTP(ctx, cl))
}
