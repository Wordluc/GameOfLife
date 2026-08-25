package core

import (
	"GameOfLife/common"
	"errors"
)

type Map[v any] struct {
	cells []*v
	size  common.Vec[int32]
}

func NewMap[v any](size common.Vec[int32]) Map[v] {
	return Map[v]{
		cells: make([]*v, size.X*size.Y),
		size:  size,
	}
}

func (m *Map[v]) SetEachElement(newCell func(x, y int32) *v) {
	var i int32
	for range m.cells {
		m.cells[i] = newCell(i%m.size.X, i/m.size.X)
		i++
	}
}

func (m *Map[v]) SetRawCell(c *v, pos common.Vec[int32]) error {
	p := pos.X + (pos.Y * m.size.X)
	if int(p) >= len(m.cells) {
		return errors.New("Out of bound")
	}
	m.cells[p] = c
	return nil
}

func (m *Map[v]) GetNeighborhoodCells(pos common.Vec[int32], size common.Vec[int32]) (res map[common.Vec[int32]]*v, err error) {
	if size.X <= 0 || size.Y <= 0 {
		return nil, errors.New("invalid neighborhood size")
	}
	if size.X%2 == 0 || size.Y%2 == 0 {
		return nil, errors.New("neighborhood size must be odd")
	}
	halfX := size.X / 2
	halfY := size.Y / 2

	res = make(map[common.Vec[int32]]*v, size.X*size.Y)
	var worldX, worldY, idx int32
	for ix := range size.X {
		for iy := range size.Y {
			worldX = pos.X + (ix - halfX)
			worldY = pos.Y + (iy - halfY)

			if worldX < 0 || worldY < 0 || worldX >= m.size.X || worldY >= m.size.Y {
				continue
			}

			idx = worldX + worldY*m.size.X
			res[common.Vec[int32]{X: worldX, Y: worldY}] = m.cells[idx]
		}
	}
	return res, nil
}

func (m *Map[v]) GetCell(pos common.Vec[int32]) (res *v, err error) {
	if pos.X < 0 || pos.Y < 0 || pos.X >= m.size.X || pos.Y >= m.size.Y {
		return nil, errors.New("Out of bound")
	}
	return m.cells[pos.X+pos.Y*m.size.X], nil
}

func (m *Map[v]) ForEach(f func(x, y int32) error) (err error) {
	for iy := range m.size.Y {
		for ix := range m.size.X {
			err = f(ix, iy)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
