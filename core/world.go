package core

import (
	"GameOfLife/common"
	"errors"
	"fmt"
	"slices"
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
	w.people = append(w.people, p)
	PerformPathFindig(w, p.currentCell.position, common.Vec[int32]{X: 20, Y: 20})
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
	for pos := range neighborhood {
		cell := new(definition.Construtor(pos))
		err = w.Map.SetRawCell(cell, pos)
		if err != nil {
			return err
		}
		w.cellBlock[cellType] = append(w.cellBlock[cellType], cell)
	}
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
	//	for _, people := range peopleByJob {
	//		for _, person := range people {
	////			paths := PerformPathFindig(w, person.currentCell.position, common.Vec[int32]{X: 20, Y: 20})
	////			if len(paths) != 0 {
	////				fmt.Printf("%v\n", paths)
	////			}
	//			return nil
	//			//	xOffset, yOffset := rand.Int31n(3), rand.Int31n(3)
	//			//	cell = person.currentCell
	//			//	newP := common.Vec[int32]{X: cell.position.X + xOffset - 1, Y: cell.position.Y + yOffset - 1}
	//			//	newCell, err := w.Map.GetCell(newP)
	//			//	if err != nil {
	//			//		continue
	//			//	}
	//			//	t := cell.PopPerson(func(p Person) bool { return p.id == person.id })
	//			//	if t != nil {
	//			//		person.Touch()
	//			//		newCell.AppendPerson(t)
	//			//	}
	//
	//		}
	//
	//	}
	return nil
}
