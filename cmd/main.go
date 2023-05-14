package main

import (
	"context"
	"fmt"

	flavGenerator "github.com/ilbarlo/flavourGeneratorProducer/pkg/flavourgenerator"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

func main() {
	ctx := context.Background()

	cl, err := flavGenerator.GetKClient(ctx)
	utilruntime.Must(err)

	flavGenerator.StartController(cl)

	fmt.Println("Started controller")
	select {}
}
