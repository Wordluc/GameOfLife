package core

import "GameOfLife/common"

type Job string

const (
	FARMER     = "FARMER"
	MINER      = "MINER"
	WOODCUTTER = "WOODCUTTER"
	FISHERMAN  = "FISHERMAN"
)

var JobToCell = map[Job][]CellType{
	FARMER:     {WHEAT_FIELD},
	MINER:      {MINE},
	WOODCUTTER: {FOREST},
	FISHERMAN:  {DOCK},
}

type Person struct {
	id          int
	Job         Job
	currentCell *BaseCell
	touch       int
	paths       *common.Queue[common.Vec[int32]]
	isWorking   bool
}

func newPerson(job Job) Person {
	p := Person{
		id:  ID_PEOPLE,
		Job: job,
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
