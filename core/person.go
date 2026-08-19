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
	STOP    = "STOP"
)

type Person struct {
	id          int
	idOrigin    int
	Job         Job
	currentCell *BaseCell
	touch       int
	paths       *common.Queue[common.Vec[int32]]
	Status      Status
}

func newPerson(job Job, idOrigin int) Person {
	p := Person{
		id:       ID_PEOPLE,
		Job:      job,
		idOrigin: idOrigin,
	}
	ID_PEOPLE++
	return p
}

func (p *Person) Touch() {
	p.touch = TOUCH_ID
}

func (p *Person) IsTouch() bool {
	return p.touch == TOUCH_ID
}
