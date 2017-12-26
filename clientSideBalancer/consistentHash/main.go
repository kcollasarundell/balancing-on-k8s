package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/kcollasarundell/balancing-on-k8s/clientSideBalancer/base"
	"github.com/kcollasarundell/balancing-on-k8s/clientSideBalancer/base/consistentHashBalancer"
	"github.com/kcollasarundell/balancing-on-k8s/rng"
	"go.uber.org/ratelimit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
)

type contextValue string

const (
	address = "dns:///rng-headless"
	key     = contextValue("userid")
)

var userids = []string{"beth", "bob", "barb", "barry", "beatrice", "brian"}

func main() {
	ctx := base.Core()
	b := consistentHashBalancer.NewBuilder(key)
	balancer.Register(b)

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithBalancerName(b.Name()), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := rng.NewRngClient(conn) // create an rng client on the above grpc connection
	rl := ratelimit.New(10, ratelimit.WithoutSlack)
	// Let's randomly select a user to make the request as
	// (todo add some logging or metrics around userid in the backend)
	rand.Seed(time.Now().Unix())

	// Make lots of requests
	for {
		ctx := context.WithValue(ctx, key, userids[rand.Int63n(int64(len(userids)))])

		if ctx.Err() != nil {
			log.Fatalf("What:%s ", ctx.Err())
			break
		}
		rl.Take()

		r, err := c.Rng(ctx, &rng.Source{Name: "consistentHash"})
		if err != nil {
			log.Printf("Wat: %v", err)
		}
		log.Printf("request %v, %d", ctx.Value(key), r.GetRN())

	}

}
