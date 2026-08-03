package core

type CellType int

const (
	GRASS = iota
	STONE
	HOUSE
)

type ICell interface {
	GetType() CellType
	SetType(CellType)
	GetPeopleNumber() int32
	SetPeopleNumber(int32)
	Touch()
	IsTouch() bool
}

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
	b.people = n
}

func NewGrassCell() BaseCell {
	return BaseCell{
		blockType: GRASS,
	}
}

func NewStoneCell() BaseCell {
	return BaseCell{
		blockType: STONE,
	}
}
