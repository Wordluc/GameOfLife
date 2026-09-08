package core

import (
	"GameOfLife/common"
	"errors"
)

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
	Id       ID_AGENT
	IdNation ID_NATION
	Job      Job
	touch    int
	paths    *common.Queue[common.Vec[int32]]
	Status   AgentStatus
	pos      common.Vec[int32]
}

func newAgent(job Job, idNation ID_NATION, pos common.Vec[int32]) Agent {
	p := Agent{
		Id:       CURRENT_ID_AGENT,
		Job:      job,
		IdNation: idNation,
		pos:      pos,
	}
	CURRENT_ID_AGENT++
	return p
}

func (a *Agent) FollowPath_StarA() (to *common.Vec[int32], err error) {
	if a.paths == nil {
		return nil, nil
	}
	from := a.paths.GetBack(1)
	if from != nil && !from.IsEqual(a.pos) {
		return nil, errors.New("Error initial position")
	}
	{
		to, end := a.paths.Denqueue()
		if end {
			a.Status = WORKING
			return nil, nil
		}
		a.Status = MOVING
		a.pos = to
	}
	return &a.pos, nil
}
