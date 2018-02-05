package main

import (
	"time"

	"google.golang.org/grpc"

	"google.golang.org/grpc/naming"
)

func main() {
	r, _ := naming.NewDNSResolverWithFreq(time.Minute)
	b := grpc.RoundRobin(r)
	grpc.Dial("serviceName", grpc.WithBalancer(b))

	defer d()
}
