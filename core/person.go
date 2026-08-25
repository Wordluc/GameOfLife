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

type Agent struct {
	id       ID_AGENT
	idNation ID_NATION
	Job      Job
	touch    int
	paths    *common.Queue[common.Vec[int32]]
	Status   AgentStatus
	pos      common.Vec[int32]
}

func newPerson(job Job, idNation ID_NATION, pos common.Vec[int32]) Agent {
	p := Agent{
		id:       CURRENT_ID_AGENT,
		Job:      job,
		idNation: idNation,
		pos:      pos,
	}
	CURRENT_ID_AGENT++
	return p
}

func (p *Agent) TouchMOVE() {
	p.touch = TOUCH_MOVE_PERSON_ID
}

func (p *Agent) IsTouchMOVE() bool {
	return p.touch == TOUCH_MOVE_PERSON_ID
}
