package core

import (
	"GameOfLife/common"
	"errors"
	"fmt"
	"slices"
	"sync"
)

var TOUCH_ID int
var ID_PEOPLE int = 1

type World struct {
	Map       Map
	cellBlock map[CellType][]*BaseCell
	people    []*Person
	resources map[string]int
}

func NewWorld(size common.Vec[int32]) (w World) {
	w.Map = NewMap(size)
	w.cellBlock = map[CellType][]*BaseCell{}
	return w
}

func (w *World) GenerateMap() {
	w.Map.FeedMap(func(x, y int32) *BaseCell {
		return new(NewGrassCell(common.Vec[int32]{X: x, Y: y}))
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

func (w *World) NewPerson(job Job, cell *BaseCell) *Person {
	p := new(newPerson(job))
	cell.AppendPerson(p)
	p.paths = new(common.NewQueue(PerformPathFindig(w, p.currentCell.position, common.Vec[int32]{X: 20, Y: 20})))
	w.people = append(w.people, p)
	return p
}

func (w *World) followPath(person *Person) error {
	from := person.paths.GetBack(1)
	if from != nil && !from.IsEqual(person.currentCell.position) {
		return errors.New("Error initial position")
	}
	newPos, end := person.paths.Denqueue()
	if end {
		return nil
	}
	person.currentCell.PopPerson(func(p Person) bool {
		return p.id == person.id
	})
	newCell, err := w.Map.GetCell(newPos)
	if err != nil {
		return err
	}
	return newCell.AppendPerson(person)
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
	for pos := range neighborhood {
		cell := new(definition.Constructor(pos))
		err = w.Map.SetRawCell(cell, pos)
		if err != nil {
			return err
		}
		w.cellBlock[cellType] = append(w.cellBlock[cellType], cell)
	}

	var wait sync.WaitGroup
	for i := range w.people {
		wait.Go(func() {
			w.people[i].paths = new(common.NewQueue(PerformPathFindig(w, w.people[i].currentCell.position, common.Vec[int32]{X: 20, Y: 20})))
		})
	}
	wait.Wait()
	return nil
}

func (w *World) MovementSimulation() error {
	var peopleByJob map[Job][]*Person = make(map[Job][]*Person)
	for i, person := range w.people {
		if person.IsTouch() {
			continue
		}
		if _, ok := peopleByJob[person.Job]; !ok {
			peopleByJob[person.Job] = []*Person{w.people[i]}
		} else {
			peopleByJob[person.Job] = append(peopleByJob[person.Job], w.people[i])
		}
	}
	var err error
	for _, people := range peopleByJob {
		for _, person := range people {
			err = w.followPath(person)
			if err != nil {
				return err
			}
		}

	}
	return nil
}
