package resolver

import (
	"strconv"
	"time"

	"google.golang.org/grpc/resolver"
)

func NewResolver() *Resolver {
	return &Resolver{
		scheme: "steve",
	}
}

type Resolver struct {
	scheme         string
	cc             resolver.ClientConn
	bootstrapAddrs []resolver.Address
}

func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	r.cc = cc

	return r, nil
}

func (r *Resolver) Scheme() string {
	return r.scheme
}

func (*Resolver) ResolveNow(o resolver.ResolveNowOption) {

}

// Close is a noop for Resolver.
func (*Resolver) Close() {}

// NewAddress calls cc.NewAddress.
func (r *Resolver) NewAddress(addrs []resolver.Address) {
	r.cc.NewAddress(addrs)
}

// NewServiceConfig calls cc.NewServiceConfig.
func (r *Resolver) NewServiceConfig(sc string) {
	r.cc.NewServiceConfig(sc)
}

// GenerateAndRegisterManualResolver generates a random scheme and a Resolver
// with it. It also regieter this Resolver.
// It returns the Resolver and a cleanup function to unregister it.
func GenerateAndRegisterManualResolver() (*Resolver, func()) {
	scheme := strconv.FormatInt(time.Now().UnixNano(), 36)
	r := NewBuilderWithScheme(scheme)
	resolver.Register(r)
	return r, func() { resolver.UnregisterForTesting(scheme) }
}
