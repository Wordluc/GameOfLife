package core

import (
	"GameOfLife/common"
	"errors"
)

type Map struct {
	cells []*BaseCell
	size  common.Vec[int32]
}

func NewMap(size common.Vec[int32]) Map {
	return Map{
		cells: make([]*BaseCell, size.X*size.Y),
		size:  size,
	}
}

func (m *Map) ForeachCell(newCell func(x, y int32) *BaseCell) {
	var i int32
	for range m.cells {
		m.cells[i] = newCell(i%m.size.X, i/m.size.X)
		i++
	}
}

func (m *Map) SetRawCell(c *BaseCell, pos common.Vec[int32]) error {
	p := pos.X + (pos.Y * m.size.X)
	c.pos = common.Vec[int32]{X: pos.X, Y: pos.Y}
	m.cells[p] = c
	return nil
}

func (m *Map) GetNeighborhoodCells(pos common.Vec[int32], size common.Vec[int32]) (res map[common.Vec[int32]]*BaseCell, err error) {
	if size.X <= 0 || size.Y <= 0 {
		return nil, errors.New("invalid neighborhood size")
	}
	if size.X%2 == 0 || size.Y%2 == 0 {
		return nil, errors.New("neighborhood size must be odd")
	}
	halfX := size.X / 2
	halfY := size.Y / 2

	res = make(map[common.Vec[int32]]*BaseCell, size.X*size.Y)
	for iy := range size.Y {
		for ix := range size.X {
			worldX := pos.X + (ix - halfX)
			worldY := pos.Y + (iy - halfY)

			if worldX < 0 || worldY < 0 || worldX >= m.size.X || worldY >= m.size.Y {
				continue
			}

			idx := worldX + worldY*m.size.X
			pos := common.Vec[int32]{X: worldX, Y: worldY}
			res[pos] = m.cells[idx]
		}
	}
	return res, nil
}

func (m *Map) GetCell(pos common.Vec[int32]) (res *BaseCell, err error) {
	if pos.X < 0 || pos.Y < 0 || pos.X >= m.size.X || pos.Y >= m.size.Y {
		return nil, errors.New("Out of bound")
	}
	return m.cells[pos.X+pos.Y*m.size.X], nil
}

func (m *Map) ForEach(f func(x, y int32, cell *BaseCell) error) (err error) {
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
