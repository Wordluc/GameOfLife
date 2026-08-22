package core

import (
	"GameOfLife/common"
	"cmp"
	"errors"
	"fmt"
	"math/rand"
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
	Populations        map[ID_NATION]*WorldPopulation
	resources          map[ID_NATION]map[Resource]float32
	toRunPathFinding   bool
	CellType_ToPosCell bindingCellTypeToCell
	IdNations          []ID_NATION
}

func NewWorld(size common.Vec[int32]) (w *World) {
	w = &World{}
	w.Map = new(NewMap(size))
	w.resources = map[ID_NATION]map[Resource]float32{}
	w.Populations = map[ID_NATION]*WorldPopulation{}
	w.CellType_ToPosCell = make(bindingCellTypeToCell)
	return w
}

func (w *World) GenerateMap() {
	var c *BaseCell
	w.Map.ForeachCell(func(x, y int32) *BaseCell {
		c = NewEmptyBaseCell(common.Vec[int32]{X: x, Y: y})
		w.CellType_ToPosCell.SetCellTypeCell(c, GRASS)
		return c
	})
}

type From = common.Vec[int32]
type To = common.Vec[int32]

var cachedPath map[From]map[To][]common.Vec[int32] = map[From]map[To][]common.Vec[int32]{}

func (w *World) setNewPathFinding(people ...*Person) {
	var cellTypeToGo CellType
	var posCellsCouldGo []common.Vec[int32]
	var person *Person
	var path []common.Vec[int32]
	var cell *BaseCell
	var goal common.Vec[int32]
	var end bool
	//GATHER POSSIBLE PATHS, SORTED BY DISTANCE
	for _, person = range people {
		if person.status == DEAD {
			continue
		}
		cellTypeToGo = JobToCells[person.Job]
		posCellsCouldGo = w.CellType_ToPosCell[cellTypeToGo]
		goals := common.NewQueue(
			posCellsCouldGo,
			func(a, b common.Vec[int32]) int {
				return cmp.Compare(fromAtoB(person.pos, a), fromAtoB(person.pos, b))
			})
		if cachedPath[person.pos] == nil {
			cachedPath[person.pos] = make(map[To][]common.Vec[int32])
		}
		for {
			goal, end = goals.Denqueue()
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
				//REMOVE FIRST ELEMENT, THE ORIGIN (person.pos)
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
	for _, n := range neighborhood {
		if definition.ConvertFrom != nil {
			if !slices.Contains(definition.ConvertFrom, n.GetType()) {
				return fmt.Errorf("%v not support in %v", cellType, neighborhood[pos].GetType())
			}
		}
		if definition.CanConvert != nil {
			if !definition.CanConvert(pos, w.Map) {
				return fmt.Errorf("CanConvert %v failed", cellType)
			}
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
	for i := range w.IdNations {
		w.setNewPathFinding(w.Populations[w.IdNations[i]].people.GetAll()...)
	}
	w.toRunPathFinding = false
}

func (w *World) NewNation(idNation ID_NATION, resource map[Resource]float32) {
	w.Populations[idNation] = new(NewWorldPopulation(w))
	w.IdNations = append(w.IdNations, idNation)
	w.resources[idNation] = resource
}

func (w *World) NewPerson(job Job, where common.Vec[int32], idNation ID_NATION) *Person {
	if !slices.Contains(w.IdNations, idNation) {
		w.NewNation(idNation, map[Resource]float32{FOOD: 1000})
	}
	p := w.Populations[idNation].newPerson(job, where, idNation)
	w.setNewPathFinding(p)
	p.TouchMOVE()
	return p
}

func (w *World) MovementSimulation() (err error) {
	for i := range w.IdNations {
		err = w.Populations[w.IdNations[i]].movePopulationToGoals(w.IdNations[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *World) HarvestingSimulation() error {
	for _, idNation := range w.IdNations {
		for celltype, cellsPos := range w.CellType_ToPosCell {
			for _, pos := range cellsPos {
				n := len(w.Populations[idNation].Pos_IdNation_ToIdPeople[pos][idNation])
				for _, q := range CellTypeToResource[celltype] {
					w.resources[idNation][q.What] += float32(n) * q.Amount
				}
			}
		}

		for _, person := range w.Populations[idNation].people.GetAll() {
			if person.status == DEAD {
				continue
			}
			for _, q := range JobToConsumingCost[person.Job] {
				w.resources[person.idNation][q.What] -= q.Amount
			}

		}
	}
	return nil
}

func (w *World) StarvingSimulation() error {
	for _, idNation := range w.IdNations {
		if w.resources[idNation][FOOD] > -100 {
			continue
		}
		var maxTime = 10
		population := w.Populations[idNation].people.GetAll()
		r := rand.Intn(len(population))
		var person *Person
		for {
			if maxTime == 0 {
				continue
			}
			person = population[r]
			//TO OPTIMIZE
			if person.status == DEAD {
				r = rand.Intn(len(population))
				maxTime--
				break
			}
			if person.paths != nil {
				lastPosCell := person.paths.GetLast()
				lastCell, _ := w.Map.GetCell(*lastPosCell)
				lastCell.VirtualNPopulation--
			}
			person.status = DEAD
			w.Populations[idNation].Pos_IdNation_ToIdPeople.RemovePerson(person)
			w.toRunPathFinding = true
			break
		}
	}
	return nil
}
