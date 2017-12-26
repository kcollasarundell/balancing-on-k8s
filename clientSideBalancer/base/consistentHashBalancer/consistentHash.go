// Package consistentHashBalancer is all about building balancers that use a consistent hash
// keyed off a specific context.Value.
package consistentHashBalancer

import (
	"fmt"
	"sync"

	"golang.org/x/net/context"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/resolver"
)

// Name is the name of consistentHash balancer.
const baseName = "consistentHash"

// NewBuilder creates a new consistentHash balancer builder
// that keys off of a context value for field.
// it requires a two step calling process for each value used.
// b := consistentHash.NewBuilder("keyName")
// balancer.Register(b)
// grpc.WithBalancerName(b.Name)
func NewBuilder(key interface{}) balancer.Builder {
	name := baseName + fmt.Sprintf("%v", key)
	return base.NewBalancerBuilder(name, &hashBalancerBuilder{Name: name, key: key})
}

type hashBalancerBuilder struct {
	Name string
	key  interface{}
}

func (h *hashBalancerBuilder) Build(readySCs map[resolver.Address]balancer.SubConn) balancer.Picker {
	grpclog.Infof("consistentHashBuilder: newPicker called with key: %v and readySCs: %v", h.key, readySCs)
	r := newRing(100)

	for address, sc := range readySCs {
		r.addNode(address, sc)

	}
	return &hashPicker{
		key:      h.key,
		subConns: r,
	}
}

type hashPicker struct {
	// subConns is the ring of subConns when this picker was
	// created. The ring is immutable. Each Get() will do a get on the ring
	// to find the appropriate backend based on the context.Value for key
	subConns *ring
	// key is used to find the appropriate field in the context to hash
	key interface{}

	mu sync.RWMutex
}

func (p *hashPicker) Pick(ctx context.Context, opts balancer.PickOptions) (balancer.SubConn, func(balancer.DoneInfo), error) {
	if p.subConns.size <= 0 {
		return nil, nil, balancer.ErrNoSubConnAvailable
	}
	value := ctx.Value(p.key)
	s, ok := value.(string)
	if !ok {
		return nil, nil, fmt.Errorf("Invalid request, context does not contain value %v", p.key)
	}

	sc := p.subConns.get(s)
	return *sc.subConn, nil, nil
}
