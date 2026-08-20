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
	cellType           CellType
	maxNPopulation     int
	VirtualNPopulation int
	other              any
	touch              int
	pos                common.Vec[int32]
}

func NewEmptyBaseCell(pos common.Vec[int32]) *BaseCell {
	return &BaseCell{
		pos:      pos,
		cellType: GRASS,
	}
}

func (BaseCell *BaseCell) Touch() {
	BaseCell.touch = TOUCH_ID
}

func (BaseCell *BaseCell) IsTouch() bool {
	return BaseCell.touch == TOUCH_ID
}

func (BaseCell *BaseCell) SetType(t CellType) {
	BaseCell.cellType = t
}

func (BaseCell *BaseCell) GetType() CellType {
	return BaseCell.cellType
}

func (BaseCell *BaseCell) GetPos() common.Vec[int32] {
	return BaseCell.pos
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
			for c, typeCell := range neir {
				diff := c.Sub(pos)
				if diff.X != 0 && diff.Y != 0 {
					continue
				}
				if slices.ContainsFunc([]CellType{DOCK, GRASS}, func(a CellType) bool {
					return a == typeCell.cellType
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
	c.cellType = GRASS
	c.touch = TOUCH_ID
	return nil
}

func ConvertToStoneCell(c *BaseCell) error {
	c.cellType = STONE
	c.touch = TOUCH_ID
	return nil
}

func ConvertToHouseCell(c *BaseCell) error {
	c.cellType = HOUSE
	c.touch = TOUCH_ID
	return nil
}
func ConvertToWheatCell(c *BaseCell) error {
	c.cellType = WHEAT_FIELD
	c.touch = TOUCH_ID
	c.maxNPopulation = 10
	return nil
}
func ConvertToMineCell(c *BaseCell) error {
	c.cellType = MINE
	c.touch = TOUCH_ID
	c.maxNPopulation = 10
	return nil
}
func ConvertToForestCell(c *BaseCell) error {
	c.cellType = FOREST
	c.touch = TOUCH_ID
	c.maxNPopulation = 10
	return nil
}
func ConvertToWaterCell(c *BaseCell) error {
	c.cellType = WATER
	c.touch = TOUCH_ID
	return nil
}
func ConvertToDockCell(c *BaseCell) error {
	c.cellType = DOCK
	c.touch = TOUCH_ID
	c.maxNPopulation = 10
	return nil
}
