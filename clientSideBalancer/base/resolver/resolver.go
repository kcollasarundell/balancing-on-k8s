package resolver

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/kcollasarundell/balancing-on-k8s/resolve"
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
	ctx       context.Context
}

func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	r.cc = cc
	address := target.Authority

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(10*time.Second))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	r.source = resolve.NewResolveClient(conn)
	c, err := r.source.ResolveStream(context.Background(), &resolve.Source{Name: target.Endpoint})

	go func() {
		for {
			rawAddresses, err := c.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("%v.resolveStream(_) = _, %v", r.source, err)
			}

			var backends []resolver.Address
			for _, rawAddress := range rawAddresses.Name {
				backends = append(backends, resolver.Address{
					Addr:       rawAddress,
					ServerName: rawAddress,
				})
			}
			r.cc.NewAddress(backends)

		}
	}()

	return r, nil
}

func (r *Resolver) Scheme() string {
	return r.scheme
}

func (*Resolver) ResolveNow(o resolver.ResolveNowOption) {}

func (*Resolver) Close() {}

// func (r *Resolver) NewAddress(addrs []resolver.Address) {
// 	r.cc.NewAddress(addrs)
// }

// func (r *Resolver) NewServiceConfig(sc string) {
// 	r.cc.NewServiceConfig(sc)
// }
