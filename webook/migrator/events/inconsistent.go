package events

import "context"

type Producer interface {
	ProducerInconsistentEvent(ctx context.Context, evt InconsistentEvent) error
}

const (
	// InconsistentEventyTypeNEQ 不一致
	InconsistentEventyTypeNEQ = "neq"
	// InconsistentEventyTypeTargetMissing  校验的目标数据缺了这一条
	InconsistentEventTypeTargetMissing = "target_missing"
	//InconsistentEventyTypeBaseMissing  校验的源数据缺了这一条
	InconsistentEventTypeBaseMissing = "base_missing"
)

type InconsistentEvent struct {
	ID int64
	//取值为 src 以源表为准, 取值为dst 以目标为准
	Direction string
	//有些时候,一些观测 第三方需要知道是什么引起的不一致
	Type string
}
