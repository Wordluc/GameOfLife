package core

import "GameOfLife/common"

type Job string

const (
	FARMER     = "FARMER"
	MINER      = "MINER"
	WOODCUTTER = "WOODCUTTER"
	FISHERMAN  = "FISHERMAN"
)

type Status string

const (
	WORKING = "WORKING"
	MOVING  = "MOVING"
	DEAD    = "DEAD"
	IDLE    = "IDLE"
	STOP    = "STOP"
)

type Person struct {
	id       ID_PERSON
	idNation ID_NATION
	Job      Job
	touch    int
	paths    *common.Queue[common.Vec[int32]]
	status   Status
	pos      common.Vec[int32]
}

func newPerson(job Job, idNation ID_NATION, pos common.Vec[int32]) Person {
	p := Person{
		id:       CURRENT_ID_PEOPLE,
		Job:      job,
		idNation: idNation,
		pos:      pos,
	}
	CURRENT_ID_PEOPLE++
	return p
}

func (p *Person) TouchMOVE() {
	p.touch = TOUCH_MOVE_PERSON_ID
}

func (p *Person) IsTouchMOVE() bool {
	return p.touch == TOUCH_MOVE_PERSON_ID
}
