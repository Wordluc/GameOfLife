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
)

type AgentStatus string

const (
	WORKING AgentStatus = "WORKING"
	MOVING  AgentStatus = "MOVING"
	DEAD    AgentStatus = "DEAD"
	IDLE    AgentStatus = "IDLE"
	STOP    AgentStatus = "STOP"
)

type Agent interface {
	Id() ID_AGENT
	Status() AgentStatus
	Pos() common.Vec[int32]
}

type agentCore struct {
	id     ID_AGENT
	status AgentStatus
	pos    common.Vec[int32]
}

func newAgentCore(pos common.Vec[int32]) agentCore {
	a := agentCore{
		id:  CURRENT_ID_AGENT,
		pos: pos,
	}
	CURRENT_ID_AGENT++
	return a
}

func (a *agentCore) Id() ID_AGENT {
	return a.id
}

func (a *agentCore) SetId(id ID_AGENT) {
	a.id = id
}

func (a *agentCore) Status() AgentStatus {
	return a.status
}

func (a *agentCore) SetStatus(status AgentStatus) {
	a.status = status
}

func (a *agentCore) Pos() common.Vec[int32] {
	return a.pos
}

func (a *agentCore) SetPos(pos common.Vec[int32]) {
	a.pos = pos
}

type pathFollower struct {
	agentCore
	paths *common.Queue[common.Vec[int32]]
}

func (a *pathFollower) Paths() *common.Queue[common.Vec[int32]] {
	return a.paths
}

func (a *pathFollower) SetPaths(paths *common.Queue[common.Vec[int32]]) {
	a.paths = paths
}

func (a *pathFollower) FollowPath_StarA() (to *common.Vec[int32], err error) {
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
			a.status = WORKING
			return nil, nil
		}
		a.status = MOVING
		a.pos = to
	}
	return &a.pos, nil
}

type Pawn struct {
	pathFollower
	idNation ID_NATION
	job      Job
}

func NewPawn(job Job, idNation ID_NATION, pos common.Vec[int32]) *Pawn {
	return &Pawn{
		pathFollower: pathFollower{agentCore: newAgentCore(pos)},
		idNation:     idNation,
		job:          job,
	}
}

func (p *Pawn) IdNation() ID_NATION {
	return p.idNation
}

func (p *Pawn) SetIdNation(idNation ID_NATION) {
	p.idNation = idNation
}

func (p *Pawn) Job() Job {
	return p.job
}

func (p *Pawn) SetJob(job Job) {
	p.job = job
}
