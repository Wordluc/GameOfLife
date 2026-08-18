package core

import (
	"GameOfLife/common"
	"errors"
	"fmt"
	"math"
	"slices"
)

var TOUCH_ID int
var ID_PEOPLE int = 1

type World struct {
	Map              *Map
	cellBlock        map[CellType]map[common.Vec[int32]]*BaseCell
	people           []*Person
	resources        map[Resource]float32
	toRunPathFinding bool
}

func NewWorld(size common.Vec[int32]) (w World) {
	w.Map = new(NewMap(size))
	w.cellBlock = map[CellType]map[common.Vec[int32]]*BaseCell{}
	w.resources = map[Resource]float32{}
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
	var cachePath map[common.Vec[int32]]map[common.Vec[int32]]*common.Queue[common.Vec[int32]] = map[common.Vec[int32]]map[common.Vec[int32]]*common.Queue[common.Vec[int32]]{}
	getMinPath := func(person *Person, cells map[common.Vec[int32]]*BaseCell) (goal *BaseCell, pathToGoal *common.Queue[common.Vec[int32]], found bool) {
		var path *common.Queue[common.Vec[int32]]
		var minPath int = math.MaxInt
		var oldGoalPos common.Vec[int32]
		if person.paths != nil {
			oldGoal, _ := w.Map.GetCell(*person.paths.GetLast())
			oldGoal.VirtualNPopulation--
			oldGoalPos = oldGoal.pos
			if oldGoal.VirtualNPopulation < 0 {
				oldGoal.VirtualNPopulation = 0
			}
		}
		for pos := range cells {
			if cells[pos].maxNPopulation != 0 && cells[pos].maxNPopulation < cells[pos].VirtualNPopulation+1 {
				continue
			}
			if cachePath[person.currentCell.pos] != nil && cachePath[person.currentCell.pos][pos] != nil {
				path = cachePath[person.currentCell.pos][pos].Clone()
			} else {
				path = common.NewQueue(PerformPathFindig(w.Map, person.currentCell.pos, pos))
				if cachePath[person.currentCell.pos] == nil {
					cachePath[person.currentCell.pos] = map[common.Vec[int32]]*common.Queue[common.Vec[int32]]{}
				}
				cachePath[person.currentCell.pos][pos] = path.Clone()
			}
			if minPath > path.Len() {
				minPath = path.Len()
				pathToGoal = path
				goal = cells[pos]
			}
		}
		if goal == nil {
			return nil, nil, false
		}
		goal.VirtualNPopulation++
		if goal.pos.IsEqual(oldGoalPos) {
			return nil, nil, true

		}
		person.paths = pathToGoal
		return goal, pathToGoal, true
	}
	for iPerson := range person {
		var i int
		found := false
		toGo := JobToCell[person[iPerson].Job]
		for {
			cellsToGo := w.cellBlock[toGo[i]]
			_, _, found = getMinPath(person[iPerson], cellsToGo)
			if found {
				break
			}
			i++
			if i >= len(toGo) {
				break
			}
		}
		if !found {
			person[iPerson].paths = nil
		}

	}
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
	person.Touch()
	return newCell.AppendPerson(person)
}

func (w *World) GetCellBlock(cellType CellType) (res []*BaseCell) {
	for _, r := range w.cellBlock[cellType] {
		res = append(res, r)
	}
	return res
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
				return fmt.Errorf("%v not support in %v", cellType, neighborhood[pos].GetType())
			}
		}
	}
	if definition.ExtraCheck != nil {
		if !definition.ExtraCheck(pos, w.Map) {
			return fmt.Errorf("ExtraCheck %v failed", cellType)
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
	w.toRunPathFinding = true
	return nil
}

func (w *World) PerformPathFinding() {
	if !w.toRunPathFinding {
		return
	}
	w.setWhereToGo(w.people...)
	w.toRunPathFinding = false
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

func (w *World) ResourcesCounting() error {
	countPeople := func(cells map[common.Vec[int32]]*BaseCell) (r int) {
		for i := range cells {
			r += cells[i].GetPeopleNumber()
		}
		return r
	}
	for celltype, cells := range w.cellBlock {
		for _, r := range CellToResource[celltype] {
			w.resources[r.What] += float32(countPeople(cells)) * r.Amount
		}
	}
	return nil
}
