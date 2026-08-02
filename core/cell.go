package core

type CellType int

const (
	GRASS = iota
	STONE
)

type ICell interface {
	GetType() CellType
	GetPeopleNumber() int
	SetPeopleNumber(int)
}

type BaseCell struct {
	blockType CellType
	people    int
	other     any
}

func (b *BaseCell) GetType() CellType {
	return b.blockType
}

func (b *BaseCell) GetPeopleNumber() int {
	return b.people
}

func (b *BaseCell) SetPeopleNumber(n int) {
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
