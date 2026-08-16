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
	WHEAT
	MINE
	DEBUG
)

type BaseCell struct {
	blockType          CellType
	people             []*Person
	maxNPopulation     int
	VirtualNPopulation int
	other              any
	touch              int
	pos                common.Vec[int32]
}

func NewEmptyBaseCell(pos common.Vec[int32]) *BaseCell {
	return &BaseCell{
		pos: pos,
	}
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

func (BaseCell *BaseCell) PopPerson(id int) (p *Person, err error) {
	for i := range BaseCell.people {
		if BaseCell.people[i].id == id {
			p = BaseCell.people[i]
			BaseCell.people[i].currentCell = nil
			BaseCell.people = slices.Delete(BaseCell.people, i, i+1)
			return p, nil
		}
	}
	return nil, fmt.Errorf("PopPerson: Id not found %v", id)
}

func (BaseCell *BaseCell) PopPersonf(check func(p Person) bool) (p *Person) {
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
	convert  func(*BaseCell) error
	WhereCan []CellType
}

var cellsDefinition map[CellType]CellDefinition = map[CellType]CellDefinition{
	GRASS: {
		convert: ConvertToGrassCell,
	},
	STONE: {
		convert:  ConvertToStoneCell,
		WhereCan: []CellType{GRASS, STONE},
	},
	HOUSE: {
		convert:  ConvertToHouseCell,
		WhereCan: []CellType{GRASS, HOUSE},
	},
	WHEAT: {
		convert:  ConvertToWheatCell,
		WhereCan: []CellType{GRASS, WHEAT},
	},
	MINE: {
		convert:  ConvertToMineCell,
		WhereCan: []CellType{STONE, MINE},
	},
}

func GetCellDefinition(t CellType) (CellDefinition, error) {
	definition, ok := cellsDefinition[t]
	if !ok {
		return CellDefinition{}, errors.New("Cell Definition not found for " + fmt.Sprint(t))
	}
	return definition, nil
}

func ConvertToGrassCell(c *BaseCell) error {
	c.blockType = GRASS
	c.touch = TOUCH_ID
	return nil
}

func ConvertToStoneCell(c *BaseCell) error {
	c.blockType = STONE
	c.touch = TOUCH_ID
	return nil
}

func ConvertToHouseCell(c *BaseCell) error {
	c.blockType = HOUSE
	c.touch = TOUCH_ID
	return nil
}
func ConvertToWheatCell(c *BaseCell) error {
	c.blockType = WHEAT
	c.touch = TOUCH_ID
	c.maxNPopulation = 10
	return nil
}
func ConvertToMineCell(c *BaseCell) error {
	c.blockType = MINE
	c.touch = TOUCH_ID
	c.maxNPopulation = 10
	return nil
}
