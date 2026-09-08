package core

import (
	"GameOfLife/common"
	"errors"
)

type Pawn struct {
	BaseAgent
	path     *common.Queue[common.Vec[int32]]
	job      Job
	idNation ID_NATION
	status   AgentStatus
}

func (p *Pawn) Move() (to *common.Vec[int32], err error) {
	p.status = IDLE
	if p.path == nil {
		return nil, nil
	}
	from := p.path.GetBack(1)
	if from != nil && !from.IsEqual(p.pos) {
		return nil, errors.New("Error initial position")
	}
	{
		to, end := p.path.Denqueue()
		if end {
			p.status = WORKING
			return nil, nil
		}
		p.status = MOVING
		p.pos = to
	}
	return &p.pos, nil
}
