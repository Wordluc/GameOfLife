package core

import "fmt"

import "errors"

type CellType int

const (
	GRASS = iota
	STONE
	HOUSE
)

type BaseCell struct {
	blockType CellType
	people    int32
	other     any
	touch     int
}

func (b *BaseCell) Touch() {
	b.touch = TOUCH_ID
}

func (b *BaseCell) IsTouch() bool {
	return b.touch == TOUCH_ID
}

func (b *BaseCell) SetType(t CellType) {
	b.blockType = t
}

func (b *BaseCell) GetType() CellType {
	return b.blockType
}

func (b *BaseCell) GetPeopleNumber() int32 {
	return b.people
}

func (b *BaseCell) SetPeopleNumber(n int32) {
	if n < 0 {
		return
	}
	b.people = n
}

type CellDefinition struct {
	Construtor func() BaseCell
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

func NewGrassCell() BaseCell {
	return BaseCell{
		blockType: GRASS,
		touch:     TOUCH_ID,
	}
}

func NewStoneCell() BaseCell {
	return BaseCell{
		blockType: STONE,
		touch:     TOUCH_ID,
	}
}

func NewHouseCell() BaseCell {
	return BaseCell{
		blockType: HOUSE,
		touch:     TOUCH_ID,
	}
}
