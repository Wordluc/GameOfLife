package core

import (
	"GameOfLife/common"
	"errors"
	"fmt"
	"slices"
)

type CellType int

const (
	GRASS = iota
	STONE
	HOUSE
)

type BaseCell struct {
	blockType CellType
	people    []*Person
	other     any
	touch     int
	position  common.Vec[int32]
}

func (BaseCell *BaseCell) Touch() {
	BaseCell.touch = TOUCH_ID
}

func (BaseCell *BaseCell) IsTouch() bool {
	return BaseCell.touch == TOUCH_ID
}

func (BaseCell *BaseCell) SetType(t CellType) {
	BaseCell.blockType = t
}

func (BaseCell *BaseCell) GetType() CellType {
	return BaseCell.blockType
}

func (BaseCell *BaseCell) GetPeopleNumber() int {
	return len(BaseCell.people)
}

func (BaseCell *BaseCell) AppendPerson(p *Person) error {
	if p.id == 0 {
		return errors.New("Wrong Id")

	}
	p.currentCell = BaseCell
	BaseCell.people = append(BaseCell.people, p)
	return nil
}

func (BaseCell *BaseCell) GetPeople(check func(p Person) bool) (res []*Person) {
	for i := range BaseCell.people {
		if check(*BaseCell.people[i]) {
			res = append(res, BaseCell.people[i])
		}
	}
	return res
}

func (BaseCell *BaseCell) PopPerson(check func(p Person) bool) (p *Person) {
	for i := range BaseCell.people {
		if check(*BaseCell.people[i]) {
			p = BaseCell.people[i]
			BaseCell.people[i].currentCell = nil
			BaseCell.people = slices.Delete(BaseCell.people, i, i+1)
			return p
		}
	}
	return nil
}

type CellDefinition struct {
	Construtor func(common.Vec[int32]) BaseCell
	WhereCan   []CellType
}

var cellsDefinition map[CellType]CellDefinition = map[CellType]CellDefinition{
	GRASS: {
		Construtor: NewGrassCell,
	},
	STONE: {
		Construtor: NewStoneCell,
	},
	HOUSE: {
		Construtor: NewHouseCell,
		WhereCan:   []CellType{GRASS},
	},
}

func GetCellDefinition(t CellType) (CellDefinition, error) {
	definition, ok := cellsDefinition[t]
	if !ok {
		return CellDefinition{}, errors.New("Cell Definition not found for " + fmt.Sprint(t))
	}
	return definition, nil
}

func NewGrassCell(pos common.Vec[int32]) BaseCell {
	return BaseCell{
		blockType: GRASS,
		touch:     TOUCH_ID,
		position:  pos,
	}
}

func NewStoneCell(pos common.Vec[int32]) BaseCell {
	return BaseCell{
		blockType: STONE,
		touch:     TOUCH_ID,
		position:  pos,
	}
}

func NewHouseCell(pos common.Vec[int32]) BaseCell {
	return BaseCell{
		blockType: HOUSE,
		touch:     TOUCH_ID,
		position:  pos,
	}
}
