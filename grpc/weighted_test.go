package grpc

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"

	"github.com/ecodeclub/ekit/slice"
)

type Node struct {
	name          string
	weight        int
	currentWeight int
}

func (n *Node) Inovke() {
}

type Balancer struct {
	nodes []*Node
	lock  sync.Mutex
	t     *testing.T
}

// 平滑权重轮询
func (b *Balancer) Pick() *Node {
	b.lock.Lock()
	defer b.lock.Unlock()
	//总权重
	total := 0
	for _, n := range b.nodes {
		total += n.weight
	}
	for _, n := range b.nodes {
		n.currentWeight += n.weight
	}
	var target *Node
	for _, n := range b.nodes {
		if target == nil {
			target = n
		} else {
			if target.currentWeight <= n.currentWeight {
				target = n
			}
		}
	}
	b.t.Log("选中了", target)
	target.currentWeight = target.currentWeight - total
	b.t.Log("选中了当前节点的权重后 减去总权重", target.currentWeight)
	return target
}

// 随机轮询查询
func (b *Balancer) RandomPick() *Node {
	total := 15
	r := rand.Int31n(int32(total))
	for _, n := range b.nodes {
		r = r - int32(n.weight)
		if r < 0 {
			return n
		}
	}
	panic("错误")

}
func TestRandoSmootWRR(t *testing.T) {
	nodes := []*Node{
		{
			name:          "A",
			weight:        10,
			currentWeight: 10,
		}, {
			name:          "B",
			weight:        20,
			currentWeight: 20,
		}, {
			name:          "C",
			weight:        30,
			currentWeight: 30,
		},
	}
	noPtrNodes := slice.Map(nodes, func(idx int, src *Node) Node {
		return *src

	})
	b := &Balancer{
		nodes: nodes,
		t:     t,
	}
	for i := 1; i < 6; i++ {
		t.Log(fmt.Sprintf("选择第 %d 个节点,nodes为: %#v", i, noPtrNodes))
		target := b.RandomPick()
		target.Inovke()
		t.Log(fmt.Sprintf("选择后第 %d 个节点,nodes为: %#v", i, noPtrNodes))

	}

}

func TestSmoothWRR(t *testing.T) {
	nodes := []*Node{
		{
			name:          "A",
			weight:        10,
			currentWeight: 10,
		}, {
			name:          "B",
			weight:        20,
			currentWeight: 20,
		}, {
			name:          "C",
			weight:        30,
			currentWeight: 30,
		},
	}
	noPtrNodes := slice.Map(nodes, func(idx int, src *Node) Node {
		return *src

	})
	b := &Balancer{
		nodes: nodes,
		t:     t,
	}
	for i := 1; i < 6; i++ {
		t.Log(fmt.Sprintf("选择第 %d 个节点,nodes为: %#v", i, noPtrNodes))
		target := b.Pick()
		target.Inovke()
		t.Log(fmt.Sprintf("选择后第 %d 个节点,nodes为: %#v", i, noPtrNodes))

	}
}
