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
	BaseCell.touch = TOUCH_MOVE_PERSON_ID
}

func (BaseCell *BaseCell) IsTouch() bool {
	return BaseCell.touch == TOUCH_MOVE_PERSON_ID
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
	convert     func(*BaseCell) error
	ConvertFrom []CellType
	maxPeople   int
	CanConvert  func(pos common.Vec[int32], m *Map) (okToConvert bool)
}

var cellsDefinition map[CellType]CellDefinition = map[CellType]CellDefinition{
	GRASS: {
		convert: ConvertToGrassCell,
	},
	STONE: {
		convert:     ConvertToStoneCell,
		ConvertFrom: []CellType{GRASS, STONE},
	},
	HOUSE: {
		convert:     ConvertToHouseCell,
		ConvertFrom: []CellType{GRASS, HOUSE},
	},
	WHEAT_FIELD: {
		convert:     ConvertToWheatCell,
		ConvertFrom: []CellType{GRASS, WHEAT_FIELD},
		maxPeople:   10,
	},
	MINE: {
		convert:     ConvertToMineCell,
		ConvertFrom: []CellType{STONE, MINE},
		maxPeople:   10,
	},
	FOREST: {
		convert:     ConvertToForestCell,
		ConvertFrom: []CellType{FOREST, GRASS, WHEAT_FIELD},
		maxPeople:   10,
	},
	WATER: {
		convert: ConvertToWaterCell,
	},
	DOCK: {
		convert:     ConvertToDockCell,
		ConvertFrom: []CellType{WATER},
		maxPeople:   10,
		CanConvert: func(pos common.Vec[int32], m *Map) (okToConvert bool) {
			neir, _ := m.GetNeighborhoodCells(pos, common.Vec[int32]{X: 3, Y: 3})
			waterCount := 0
			dockOrGrassCount := 0
			for posNeighbord, cell := range neir {
				diff := posNeighbord.Clone().Sub(pos)
				if diff.X != 0 && diff.Y != 0 {
					// diagonal, skip
					continue
				}
				if cell.cellType == WATER && !posNeighbord.IsEqual(pos) {
					waterCount++
				}

				if slices.Contains([]CellType{DOCK, GRASS}, cell.cellType) {
					dockOrGrassCount++
				}
			}
			return dockOrGrassCount > 0 && waterCount > 0
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
	c.touch = TOUCH_MOVE_PERSON_ID
	return nil
}

func ConvertToStoneCell(c *BaseCell) error {
	c.cellType = STONE
	c.touch = TOUCH_MOVE_PERSON_ID
	return nil
}

func ConvertToHouseCell(c *BaseCell) error {
	c.cellType = HOUSE
	c.touch = TOUCH_MOVE_PERSON_ID
	return nil
}
func ConvertToWheatCell(c *BaseCell) error {
	c.cellType = WHEAT_FIELD
	c.touch = TOUCH_MOVE_PERSON_ID
	return nil
}
func ConvertToMineCell(c *BaseCell) error {
	c.cellType = MINE
	c.touch = TOUCH_MOVE_PERSON_ID
	return nil
}
func ConvertToForestCell(c *BaseCell) error {
	c.cellType = FOREST
	c.touch = TOUCH_MOVE_PERSON_ID
	return nil
}
func ConvertToWaterCell(c *BaseCell) error {
	c.cellType = WATER
	c.touch = TOUCH_MOVE_PERSON_ID
	return nil
}
func ConvertToDockCell(c *BaseCell) error {
	c.cellType = DOCK
	c.touch = TOUCH_MOVE_PERSON_ID
	return nil
}
