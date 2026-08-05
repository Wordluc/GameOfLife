package core

import (
	"GameOfLife/common"
	"errors"
	"fmt"
	"slices"
)

var TOUCH_ID int
var ID_PEOPLE int

type World struct {
	Map       Map
	cellBlock map[CellType][]*BaseCell
	people    []*Person
}

func NewWorld(size common.Vec[int32]) (w World) {
	w.Map = NewMap(size)
	w.cellBlock = map[CellType][]*BaseCell{}
	return w
}

func (w *World) GenerateMap() {
	w.Map.FeedMap(func(x, y int32) *BaseCell {
		return new(NewGrassCell())
	})
}

func (w *World) GetPeople(check func(p Person) bool) (res []*Person) {
	for i := range w.people {
		if check(*w.people[i]) {
			res = append(res, w.people[i])
		}
	}
	return res
}

func (w *World) NewPerson(job Job) *Person {
	p := new(newPerson(job))
	w.people = append(w.people, p)
	return p
}

func (w *World) GetCellBlock(cellType CellType) ([]*BaseCell, error) {
	return slices.Clone(w.cellBlock[cellType]), nil
}

func (w *World) AddBlock(cellType CellType, pos common.Vec[int32], size common.Vec[int32]) error {
	halfX := size.X / 2
	halfY := size.Y / 2
	if size.X <= 0 || size.Y <= 0 {
		return errors.New("invalid neighborhood size")
	}
	if size.X%2 == 0 || size.Y%2 == 0 {
		return errors.New("neighborhood size must be odd")
	}
	if pos.X-halfX < 0 || pos.Y-halfY < 0 || pos.X+halfX >= w.Map.size.X || pos.Y+halfY >= w.Map.size.Y {
		return errors.New("Out of bound")
	}

	neighborhood, err := w.Map.GetNeighborhoodCells(pos, size)
	if err != nil {
		return err
	}
	definition, err := GetCellDefinition(cellType)
	if err != nil {
		return err
	}
	if definition.WhereCan != nil {
		for _, n := range neighborhood {
			if !slices.Contains(definition.WhereCan, n.GetType()) {
				return fmt.Errorf("House not support in %v", neighborhood[pos].GetType())
			}
		}
	}
	for p := range neighborhood {
		cell := new(definition.Construtor())
		err = w.Map.SetRawCell(cell, p)
		if err != nil {
			return err
		}
		w.cellBlock[cellType] = append(w.cellBlock[cellType], cell)
	}
	return nil
}
