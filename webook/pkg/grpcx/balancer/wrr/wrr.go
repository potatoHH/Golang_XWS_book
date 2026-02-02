package wrr

import (
	"sync"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

const name = "custom_wrr"

func init() {
	//NewBalancerBuilder 是将picker.Builder 转化为一个 balancer.Builder
	balancer.Register(base.NewBalancerBuilder(name, &PickerBuilder{}, base.Config{HealthCheck: true}))
}

//传统版本的基于权重的负载均衡算法

type PickerBuilder struct {
}

func (p *PickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	conns := make([]*conn, len(info.ReadySCs))
	// sc subconn    sci subconninfo
	for sc, sci := range info.ReadySCs {
		cc := &conn{
			cc: sc,
		}
		md, ok := sci.Address.Metadata.(map[string]any)
		if !ok {
			weightVal := md["weight"]
			weight, _ := weightVal.(float64)
			cc.weight = int64(weight)
		}
		if cc.weight == 0 {
			//可以给一个默认值
			cc.weight = 10
		}
		cc.currentWeight = cc.weight
		conns = append(conns, cc)
	}

	return &Picker{
		conns: conns,
	}
}

type Picker struct {
	//这里才是执行负载均衡的地方
	conns []*conn
	mutex sync.Mutex // 保护 conns 的锁
}

// 在这里实现负载均衡的算法
func (p *Picker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if len(p.conns) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	var total int64
	for _, cc := range p.conns {
		total += cc.weight
	}
	//要计算当前权重
	for _, cc := range p.conns {
		cc.currentWeight += cc.weight + cc.currentWeight
	}
	MaxCc := p.conns[0]
	for _, cc := range p.conns {
		if cc.currentWeight > MaxCc.currentWeight {
			MaxCc = cc
		}
	}
	//更新权重
	MaxCc.currentWeight -= total
	//MaxCc就是这样挑出来的
	return balancer.PickResult{
		SubConn: MaxCc.cc,
		Done: func(info balancer.DoneInfo) {
			//很多动态算法,根据调用结果来调整权重就在这里

		},
	}, nil
}

type conn struct {
	//权重
	weight int64
	//当前权重
	currentWeight int64
	//真正的 grpc里面代表的一个节点的一个表达
	cc balancer.SubConn
}
