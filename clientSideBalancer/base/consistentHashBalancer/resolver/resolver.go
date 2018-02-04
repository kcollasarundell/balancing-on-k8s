package resolver

import (
	"log"
	"time"

	resolve "github.com/kcollasarundell/balancing-on-k8s/resolve-source"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

func NewResolver() *Resolver {
	return &Resolver{
		scheme: "steve",
	}
}

type Resolver struct {
	scheme    string
	cc        resolver.ClientConn
	source    resolve.ResolveClient
	knownAddr []string
}

func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	r.cc = cc
	address := "service-discovery-upstream"

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(10*time.Second))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	r.source = resolve.NewResolveClient(conn)

	return r, nil
}

func (r *Resolver) Scheme() string {
	return r.scheme
}

func (*Resolver) ResolveNow(o resolver.ResolveNowOption) {}

func (*Resolver) Close() {}

func (r *Resolver) NewAddress(addrs []resolver.Address) {
	r.cc.NewAddress(addrs)
}

func (r *Resolver) NewServiceConfig(sc string) {
	r.cc.NewServiceConfig(sc)
}
