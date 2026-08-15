package core

import (
	"GameOfLife/common"
	"errors"
	"fmt"
	"math"
	"slices"
	"sync"
)

var TOUCH_ID int
var ID_PEOPLE int = 1

type World struct {
	Map       Map
	cellBlock map[CellType]map[common.Vec[int32]]*BaseCell
	people    []*Person
	resources map[string]int
}

func NewWorld(size common.Vec[int32]) (w World) {
	w.Map = NewMap(size)
	w.cellBlock = map[CellType]map[common.Vec[int32]]*BaseCell{}
	return w
}

func (w *World) GenerateMap() {
	w.Map.ForeachCell(func(x, y int32) *BaseCell {
		return NewEmptyBaseCell(common.Vec[int32]{X: x, Y: y})
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
	w.setWhereToGo(p)
	w.people = append(w.people, p)
	return p
}

func (w *World) setWhereToGo(person ...*Person) {
	getMinDistant := func(personPos common.Vec[int32], cells map[common.Vec[int32]]*BaseCell) (cell *BaseCell) {
		var min int32 = math.MaxInt32
		cell = nil
		for pos := range cells {
			if min > common.DistanceAtoBVecShev(personPos, pos) {

				cell = cells[pos]
			}
		}
		return cell
	}
	var wait sync.WaitGroup
	for i := range person {
		wait.Go(func() {
			cellsToGo := w.cellBlock[JobToCell[person[i].Job]]
			goal := getMinDistant(person[i].currentCell.pos, cellsToGo)
			if goal == nil {
				return
			}

			person[i].paths = new(common.NewQueue(PerformPathFindig(w, person[i].currentCell.pos, goal.pos)))
		})
	}
	wait.Wait()
}

func (w *World) followPath(person *Person) error {
	if person.paths == nil {
		return nil
	}
	from := person.paths.GetBack(1)
	if from != nil && !from.IsEqual(person.currentCell.pos) {
		return errors.New("Error initial position")
	}
	newPos, end := person.paths.Denqueue()
	if end {
		return nil
	}
	_, err := person.currentCell.PopPerson(person.id)
	if err != nil {
		return err
	}
	newCell, err := w.Map.GetCell(newPos)
	if err != nil {
		return err
	}
	return newCell.AppendPerson(person)
}

func (w *World) GetCellBlock(cellType CellType) (res []*BaseCell, err error) {
	for _, r := range w.cellBlock[cellType] {
		res = append(res, r)
	}
	return res, err
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
				fmt.Printf("%v not support in %v\n", cellType, neighborhood[pos].GetType())
				return nil
			}
		}
	}
	for pos := range neighborhood {
		cell, err := w.Map.GetCell(pos)
		if err != nil {
			return err
		}
		delete(w.cellBlock[cell.blockType], pos)
		definition.convert(cell)
		if w.cellBlock[cellType] == nil {
			w.cellBlock[cellType] = make(map[common.Vec[int32]]*BaseCell)
		}
		w.cellBlock[cellType][pos] = cell
	}

	w.setWhereToGo(w.people...)
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
			if w.followPath(person) != nil {
				return err
			}
		}

	}
	return nil
}
