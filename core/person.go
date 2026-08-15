package core

import "GameOfLife/common"

type Job int8

const (
	FARMER Job = iota
	MINER
)

var JobToCell = map[Job]CellType{
	FARMER: WHEAT,
	MINER:  MINE,
}

type Person struct {
	id          int
	Job         Job
	currentCell *BaseCell
	touch       int
	paths       *common.Queue[common.Vec[int32]]
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
