package core

import (
	"GameOfLife/common"
	"cmp"
	"errors"
	"fmt"
	"slices"
)

type ID_NATION int

var TOUCH_MOVE_PERSON_ID int

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

type From = common.Vec[int32]
type To = common.Vec[int32]

var cachedPath map[From]map[To][]common.Vec[int32] = map[From]map[To][]common.Vec[int32]{}

func (w *World) setNewPathFinding(people ...*Person) {
	var cellTypeToGo CellType
	var posCellsCouldGo []common.Vec[int32]
	var whereToGo map[ID_PERSON]*common.Queue[common.Vec[int32]] = map[ID_PERSON]*common.Queue[common.Vec[int32]]{}
	var person *Person
	var path []common.Vec[int32]
	var cell *BaseCell
	//GATHER POSSIBLE PATHS, SORTED BY DISTANCE
	for _, person := range people {
		cellTypeToGo = JobToCells[person.Job]
		posCellsCouldGo = w.CellType_ToPosCell[cellTypeToGo]
		whereToGo[person.id] = common.NewQueue(
			posCellsCouldGo,
			func(a, b common.Vec[int32]) int {
				return cmp.Compare(fromAtoB(person.pos, a), fromAtoB(person.pos, b))
			})
	}
	//ASSIGN FIRST POSSIBLE PATH
	for idPerson, goals := range whereToGo {
		person = w.Populations.GetPerson(idPerson)
		if cachedPath[person.pos] == nil {
			cachedPath[person.pos] = make(map[To][]common.Vec[int32])
		}
		for {
			goal, end := goals.Denqueue()
			if end {
				person.status = IDLE
				break
			}
			if goal.IsEqual(person.pos) {
				break
			}

			if person.paths != nil {
				cell, _ = w.Map.GetCell(*person.paths.GetLast())
				cell.VirtualNPopulation--
				if cell.VirtualNPopulation < 0 {
					cell.VirtualNPopulation = 0
				}
				person.paths = nil
			}

			if cachedPath[person.pos][goal] != nil {
				path = cachedPath[person.pos][goal]
			} else {
				path = PerformPathFindig(w.Map, person.pos, goal)
				cachedPath[person.pos][goal] = path
			}

			if len(path) != 0 {
				cell, _ = w.Map.GetCell(goal)
				if cell.VirtualNPopulation >= cellsDefinition[cell.cellType].maxPeople {
					continue
				}
				cell.VirtualNPopulation++
				person.paths = common.NewQueue(path, nil)
				//REMOVE FIRST ELEMENT, THE ORIDIN (person.pos)
				person.paths.Denqueue()
				break
			}
		}
	}
}

func (w *World) AddBlock(cellType CellType, pos common.Vec[int32], size common.Vec[int32]) error {
	halfX := size.X / 2
	halfY := size.Y / 2
	if size.X <= 0 || size.Y <= 0 {
		return errors.New("Invalid neighborhood size")
	}
	if size.X%2 == 0 || size.Y%2 == 0 {
		return errors.New("Neighborhood size must be odd")
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
	if definition.ConvertFrom != nil {
		for _, n := range neighborhood {
			if !slices.Contains(definition.ConvertFrom, n.GetType()) {
				return fmt.Errorf("%v not support in %v", cellType, neighborhood[pos].GetType())
			}
		}
	}
	if definition.CanConvert != nil {
		if !definition.CanConvert(pos, w.Map) {
			return fmt.Errorf("CanConvert %v failed", cellType)
		}
	}
	for pos := range neighborhood {
		cell, err := w.Map.GetCell(pos)
		if err != nil {
			return err
		}
		w.CellType_ToPosCell.SetCellTypeCell(cell, cellType)
		err = definition.convert(cell)
		if err != nil {
			return err
		}
	}
	w.toRunPathFinding = true
	t := make(map[From]map[To][]common.Vec[int32])
	for from, paths := range cachedPath {
		for to, path := range paths {
			if slices.ContainsFunc(path, func(a common.Vec[int32]) bool { return neighborhood[a] != nil }) {
				continue
			}
			if t[from] == nil {
				t[from] = map[To][]common.Vec[int32]{}
			}
			t[from][to] = cachedPath[from][to]
		}
	}
	cachedPath = t
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
	w.setNewPathFinding(w.Populations.people.GetAll()...)
	w.toRunPathFinding = false
}

func (w *World) NewPerson(job Job, where common.Vec[int32], idNation ID_NATION) *Person {
	if !slices.Contains(w.idNations, idNation) {
		w.idNations = append(w.idNations, idNation)
	}
	p := w.Populations.newPerson(job, where, idNation)
	w.setNewPathFinding(p)
	p.TouchMOVE()
	return p
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
