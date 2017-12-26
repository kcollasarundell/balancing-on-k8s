package main

import (
	"log"

	"github.com/kcollasarundell/balancing-on-k8s/clientSideBalancer/base"
	"github.com/kcollasarundell/balancing-on-k8s/rng"
	"go.uber.org/ratelimit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
)

const (
	address = "dns:///rng-headless"
)

func main() {
	ctx := base.Core()

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithBalancerName(roundrobin.Name), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := rng.NewRngClient(conn) // create an rng client on the above grpc connection
	rl := ratelimit.New(1, ratelimit.WithoutSlack)
	// Make lots of requests
	for {
		if ctx.Err() != nil {
			log.Fatalf("What:%s ", ctx.Err())
			break
		}
		rl.Take()
		base.Request(ctx, "roundrobin", c)
	}

}
