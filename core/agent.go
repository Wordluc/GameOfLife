package core

import "GameOfLife/common"

type Job string

const (
	FARMER     = "FARMER"
	MINER      = "MINER"
	WOODCUTTER = "WOODCUTTER"
	FISHERMAN  = "FISHERMAN"
	ZOMBIE     = "ZOMBIE"
)

type AgentStatus string

const (
	WORKING AgentStatus = "WORKING"
	MOVING  AgentStatus = "MOVING"
	DEAD    AgentStatus = "DEAD"
	IDLE    AgentStatus = "IDLE"
	STOP    AgentStatus = "STOP"
)

type BaseAgent struct {
	id  int
	pos common.Vec[int32]
}

func (b *BaseAgent) GetPos() common.Vec[int32] {
	return b.pos
}

func (b *BaseAgent) SetPos(p common.Vec[int32]) {
	b.pos = p
}

func (b *BaseAgent) GetId() int {
	return b.id
}

type Agent interface {
	GetPos() common.Vec[int32]
	SetPos(common.Vec[int32])
	GetId() int
}
