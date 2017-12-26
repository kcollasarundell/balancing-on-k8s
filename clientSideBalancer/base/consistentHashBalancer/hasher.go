package consistentHashBalancer

import (
	"fmt"
	"hash/crc32"
	"sort"
	"sync"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/resolver"
)

type ring struct {
	nodes  nodes
	rounds int
	size   int
	sync.Mutex
}

func newRing(rounds int) *ring {
	return &ring{nodes: nodes{}, rounds: rounds}
}

func (r *ring) addNode(address resolver.Address, subConn balancer.SubConn) {
	r.Lock()
	defer r.Unlock()

	r.size++

	for round := 0; round < r.rounds; round++ {
		node := &node{
			id:      address.Addr,
			hashID:  hashID(address.Addr + fmt.Sprintf("%d", round)),
			subConn: &subConn,
		}
		r.nodes = append(r.nodes, node)
	}
	sort.Sort(r.nodes)
}

func (r *ring) removeNode(id string) error {
	return fmt.Errorf("Not implemented (i'm lazy and it recreates the ring for the new picker)")
}

func (r *ring) get(id string) *node {

	searchfn := func(i int) bool {
		return r.nodes[i].hashID >= hashID(id)
	}

	i := sort.Search(r.nodes.Len(), searchfn)

	if i >= r.nodes.Len() {
		i = 0
	}
	return r.nodes[i]
}

type node struct {
	id      string
	hashID  uint32
	subConn *balancer.SubConn
}

type nodes []*node

func (n nodes) Len() int           { return len(n) }
func (n nodes) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n nodes) Less(i, j int) bool { return n[i].hashID < n[j].hashID }

//----------------------------------------------------------
// Helpers
//----------------------------------------------------------

func hashID(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}
