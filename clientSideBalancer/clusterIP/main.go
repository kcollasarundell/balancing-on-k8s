package main

import (
	"log"
	"time"

	"go.uber.org/ratelimit"

	"github.com/kcollasarundell/balancing-on-k8s/clientSideBalancer/base"

	"github.com/kcollasarundell/balancing-on-k8s/rng"
	"google.golang.org/grpc"
)

const (
	address = "rng-cluster:8081"
)

func main() {
	ctx := base.Core()

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(10*time.Second))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := rng.NewRngClient(conn) // create an rng client on the above grpc connection
	rl := ratelimit.New(150)
	// Make lots of requests
	for {
		if ctx.Err() != nil {
			log.Fatalf("What:%s ", ctx.Err())
			break
		}
		rl.Take()
		base.Request(ctx, "clusterIP", c)
	}

}
