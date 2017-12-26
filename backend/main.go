package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"net/http"

	"github.com/kcollasarundell/balancing-on-k8s/rng"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":8081"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) Rng(context.Context, *rng.Empty) (*rng.RN, error) {
	requests.Inc()
	return &rng.RN{RN: rand.Int63n(100)}, nil
}

var requests = prometheus.NewCounter(prometheus.CounterOpts{Name: "randomNumbers", Help: "the count of random numbers generated"})

func init() {
	prometheus.MustRegister(requests)
}
func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	rng.RegisterRngServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":8080", nil)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
