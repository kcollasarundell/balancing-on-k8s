package base

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kcollasarundell/balancing-on-k8s/rng"
	"google.golang.org/grpc"
)

// Core sets up basic quality of life things to handle terminating all the stuff
func Core() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func(cancel context.CancelFunc) {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		cancel()
	}(cancel)

	return ctx
}

// Placeholder is just the basic structure to be redone in each of the examples (or used here)
func Placeholder(ctx context.Context, address string, grpcopts ...grpc.DialOption) rng.RngClient {

	// Set up a connection to the server.
	grpcopts = append(grpcopts, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(10*time.Second))
	conn, err := grpc.Dial(address, grpcopts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := rng.NewRngClient(conn)

	return c
}

// Request makes random number requests against the backend until told to quit
func Request(ctx context.Context, name string, c rng.RngClient) {

	r, err := c.Rng(ctx, &rng.Source{Name: name})
	if err != nil {
		log.Printf("Wat: %v", err)
	}
	log.Printf("request %d", r.GetRN())

}
