package core

import (
	"GameOfLife/common"
	"cmp"
	"errors"
	"fmt"
	"slices"
)

type ID_NATION int

var TOUCH_ID int

type bindingCellTypeToCell map[CellType][]common.Vec[int32]

func (b bindingCellTypeToCell) SetCellTypeCell(cell *BaseCell, newCellType CellType) {
	b[cell.cellType] = slices.DeleteFunc(b[cell.cellType], func(a common.Vec[int32]) bool { return a.IsEqual(cell.pos) })
	b[newCellType] = append(b[newCellType], cell.pos)
	cell.cellType = newCellType
}

type World struct {
	Map                *Map
	Populations        *WorldPopulation
	resources          map[Resource]float32
	toRunPathFinding   bool
	CellType_ToPosCell bindingCellTypeToCell
	idNations          []ID_NATION
}

func NewWorld(size common.Vec[int32]) (w *World) {
	w = &World{}
	w.Map = new(NewMap(size))
	w.resources = map[Resource]float32{
		FOOD: 1000,
	}
	w.Populations = new(NewWorldPopulation(w))
	w.CellType_ToPosCell = make(bindingCellTypeToCell)
	return w
}

func (w *World) GenerateMap() {
	w.Map.ForeachCell(func(x, y int32) *BaseCell {
		return NewEmptyBaseCell(common.Vec[int32]{X: x, Y: y})
	})
}

func (w *World) setWhereToGo(people ...*Person) {
	type From = common.Vec[int32]
	type To = common.Vec[int32]
	//	var cachedPath map[From]map[To]*common.Queue[common.Vec[int32]] = map[From]map[To]*common.Queue[common.Vec[int32]]{}
	//	getMinPath := func(person *Person, cellsPos []common.Vec[int32]) (goal *BaseCell, pathToGoal *common.Queue[common.Vec[int32]], found bool) {
	//		//		var path *common.Queue[common.Vec[int32]]
	//		//		var minPath int = math.MaxInt
	//		//		var oldGoalPos common.Vec[int32]
	//		//		if person.paths != nil {
	//		//			oldGoal, _ := w.Map.GetCell(*person.paths.GetLast())
	//		//			oldGoal.VirtualNPopulation--
	//		//			oldGoalPos = oldGoal.pos
	//		//			if oldGoal.VirtualNPopulation < 0 {
	//		//				oldGoal.VirtualNPopulation = 0
	//		//			}
	//		//		}
	//		//		for pos := range cells {
	//		//			if cells[pos].maxNPopulation != 0 && cells[pos].maxNPopulation < cells[pos].VirtualNPopulation+1 {
	//		//				continue
	//		//			}
	//		//			if cachePath[person.pos] != nil && cachePath[person.pos][pos] != nil {
	//		//				path = cachePath[person.pos][pos].Clone()
	//		//			} else {
	//		//				path = common.NewQueue(PerformPathFindig(w.Map, person.pos, pos))
	//		//				if cachePath[person.pos] == nil {
	//		//					cachePath[person.pos] = map[common.Vec[int32]]*common.Queue[common.Vec[int32]]{}
	//		//				}
	//		//				cachePath[person.pos][pos] = path.Clone()
	//		//			}
	//		//			if minPath > path.Len() {
	//		//				minPath = path.Len()
	//		//				pathToGoal = path
	//		//				goal = cells[pos]
	//		//			}
	//		//		}
	//		//		if goal == nil {
	//		//			return nil, nil, false
	//		//		}
	//		//		goal.VirtualNPopulation++
	//		//		if goal.pos.IsEqual(oldGoalPos) {
	//		//			return nil, nil, true
	//		//
	//		//		}
	//		//		person.paths = pathToGoal
	//		return goal, pathToGoal, true
	//	}
	var cellTypeToGo CellType
	var posCellsCouldGo []common.Vec[int32]
	var whereToGo map[ID_PERSON]*common.Queue[common.Vec[int32]] = map[ID_PERSON]*common.Queue[common.Vec[int32]]{}
	var person *Person
	for _, person := range people {
		cellTypeToGo = JobToCells[person.Job]
		posCellsCouldGo = w.CellType_ToPosCell[cellTypeToGo]
		whereToGo[person.id] = common.NewQueue(
			posCellsCouldGo,
			func(a, b common.Vec[int32]) int {
				return cmp.Compare(fromAtoB(person.pos, a), fromAtoB(person.pos, b))
			})
	}

	for idPerson, goals := range whereToGo {
		person = w.Populations.GetPerson(idPerson)
		for {
			goal, end := goals.Denqueue()
			if end {
				break
			}
			path := PerformPathFindig(w.Map, person.pos, goal)
			if len(path) != 0 {
				person.paths = common.NewQueue(path, nil)
				break
			}
		}
	}
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
		w.CellType_ToPosCell.SetCellTypeCell(cell, cellType)
		definition.convert(cell)
	}
	w.toRunPathFinding = true
	return nil
}

func (w *World) GetCellsByType(cellType CellType) (res []*BaseCell, err error) {
	var cell *BaseCell
	for _, pos := range w.CellType_ToPosCell[cellType] {
		cell, err = w.Map.GetCell(pos)
		if err != nil {
			return nil, err
		}
		res = append(res, cell)
	}
	return res, nil
}

func (w *World) PerformPathFinding() {
	if !w.toRunPathFinding {
		return
	}
	w.setWhereToGo(w.Populations.people...)
	w.toRunPathFinding = false
}

func (w *World) NewPerson(job Job, where common.Vec[int32], idNation ID_NATION) *Person {
	if !slices.Contains(w.idNations, idNation) {
		w.idNations = append(w.idNations, idNation)
	}
	w.toRunPathFinding = true
	return w.Populations.newPerson(job, where, idNation)
}

func (w *World) MovementSimulation() (err error) {
	for i := range w.idNations {
		err = w.Populations.movePopulationToGoals(w.idNations[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *World) ResourcesCounting() error {
	n := 0
	for _, idNation := range w.idNations {
		for celltype, cellsPos := range w.CellType_ToPosCell {
			for _, pos := range cellsPos {
				n = len(w.Populations.Pos_IdNation_ToIdPeople[pos][idNation])
				for _, q := range CellTypeToResource[celltype] {
					w.resources[q.What] += float32(n) * q.Amount
				}
			}
		}
	}
	//
	//	fmt.Printf("%v\n", w.resources)
	//	var maxTime = 10
	//	if w.resources[FOOD] < -100 {
	//		r := rand.Intn(len(w.people))
	//		for {
	//			if maxTime == 0 {
	//				return nil
	//			}
	//			if w.people[r].Status == DEAD {
	//				r = rand.Intn(len(w.people))
	//				maxTime--
	//				continue
	//			}
	//			if w.people[r].Status == MOVING && w.people[r].paths != nil {
	//				p := w.people[r].paths.GetLast()
	//				c, _ := w.Map.GetCell(*p)
	//				c.VirtualNPopulation--
	//				w.toRunPathFinding = true
	//			} else if w.people[r].Status == WORKING {
	//				w.people[r].currentCell.VirtualNPopulation--
	//			}
	//			w.people[r].Status = DEAD
	//			w.people[r].Touch()
	//		}
	//	}
	return nil
}
