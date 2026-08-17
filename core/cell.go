package core

import (
	"GameOfLife/common"
	"errors"
	"fmt"
	"slices"
)

type CellType string

const (
	GRASS       = "GRASS"
	STONE       = "STONE"
	HOUSE       = "HOUSE"
	WHEAT_FIELD = "WHEAT_FIELD"
	MINE        = "MINE"
	FOREST      = "FOREST"
	WATER       = "WATER"
	DOCK        = "DOCK"
	DEBUG       = "DEBUG"
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
		pos:       pos,
		blockType: GRASS,
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
	convert    func(*BaseCell) error
	WhereCan   []CellType
	ExtraCheck func(pos common.Vec[int32], m *Map) (okToConvert bool)
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
	WHEAT_FIELD: {
		convert:  ConvertToWheatCell,
		WhereCan: []CellType{GRASS, WHEAT_FIELD},
	},
	MINE: {
		convert:  ConvertToMineCell,
		WhereCan: []CellType{STONE, MINE},
	},
	FOREST: {
		convert:  ConvertToForestCell,
		WhereCan: []CellType{FOREST, GRASS, WHEAT_FIELD},
	},
	WATER: {
		convert: ConvertToWaterCell,
	},
	DOCK: {
		convert:  ConvertToDockCell,
		WhereCan: []CellType{WATER},
		ExtraCheck: func(pos common.Vec[int32], m *Map) (okToConvert bool) {
			neir, _ := m.GetNeighborhoodCells(pos, common.Vec[int32]{X: 3, Y: 3})
			for _, typeCell := range neir {
				if slices.ContainsFunc([]CellType{DOCK, GRASS}, func(a CellType) bool {
					return a == typeCell.blockType
				}) {
					return true
				}
			}
			return false
		},
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
	c.blockType = WHEAT_FIELD
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
func ConvertToForestCell(c *BaseCell) error {
	c.blockType = FOREST
	c.touch = TOUCH_ID
	c.maxNPopulation = 10
	return nil
}
func ConvertToWaterCell(c *BaseCell) error {
	c.blockType = WATER
	c.touch = TOUCH_ID
	return nil
}
func ConvertToDockCell(c *BaseCell) error {
	c.blockType = DOCK
	c.touch = TOUCH_ID
	c.maxNPopulation = 10
	return nil
}
