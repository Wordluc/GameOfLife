package core

import (
	"GameOfLife/common"
	"errors"
)

type Map struct {
	cells []ICell
	size  common.Vec[int32]
}

func NewMap(size common.Vec[int32]) Map {
	return Map{
		cells: make([]ICell, size.X*size.Y),
		size:  size,
	}
}

func (m *Map) FeedMap(newCell func(x, y int32) ICell) {
	var i int32
	for range m.cells {
		m.cells[i] = newCell(i%m.size.X, i/m.size.X)
		i++
	}
}

func (m *Map) SetRawCell(c ICell, pos common.Vec[int32]) error {
	p := pos.X + (pos.Y * m.size.X)
	m.cells[p] = c
	return nil
}

func (m *Map) SetCells(newCell func(x, y int32) ICell, pos, size common.Vec[int32]) error {
	for iy := range size.Y {
		for ix := range size.X {
			p := (pos.X + (ix - size.X/2) + (pos.Y+(iy-size.Y/2))*m.size.X)
			m.cells[p] = newCell(ix, iy)
		}
	}
	return nil
}

func (m *Map) GetCell(pos common.Vec[int32]) (res ICell, err error) {
	if pos.X < 0 || pos.Y < 0 {
		return nil, errors.New("Error getting cell")
	}
	return m.cells[pos.X+pos.Y*m.size.X], nil
}

func (m *Map) ForEach(f func(x, y int32, cell ICell) error) (err error) {
	for iy := range m.size.Y {
		for ix := range m.size.X {
			err = f(ix, iy, m.cells[ix+iy*m.size.X])
			if err != nil {
				return err
			}
		}
	}
	return nil
}
