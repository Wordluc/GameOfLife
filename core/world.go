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

type Player struct {
	visualPos common.Vec[float32]
	logic     common.Vec[int32]
}

type World struct {
	Map                *Map[BaseCell]
	Nation             map[ID_NATION]*Nation
	Zombies            *ZombieHorde
	resources          map[ID_NATION]map[Resource]float32
	CellType_ToPosCell bindingCellTypeToCell
	IdNations          []ID_NATION
	players            []Player

	personNeedingPathFinding *common.Queue[*Agent]
}

func NewWorld(size common.Vec[int32]) (w *World) {
	w = &World{}
	w.Map = new(NewMap[BaseCell](size))
	w.resources = map[ID_NATION]map[Resource]float32{}
	w.Nation = map[ID_NATION]*Nation{}
	w.CellType_ToPosCell = make(bindingCellTypeToCell)
	w.personNeedingPathFinding = common.NewQueue[*Agent](nil, nil)
	w.Zombies = new(NewZombieHorde(w))
	return w
}

func (w *World) GenerateMap() {
	var c *BaseCell
	w.Map.SetEachElement(func(x, y int32) *BaseCell {
		c = NewEmptyBaseCell(common.Vec[int32]{X: x, Y: y})
		w.CellType_ToPosCell.SetCellTypeCell(c, GRASS)
		return c
	})
}

func (w *World) toRunPathFinding(ps ...*Agent) {
	w.personNeedingPathFinding.Enqueue(ps...)
}

func (w *World) toRunPathFindingForAll() {
	for i := range w.IdNations {
		w.personNeedingPathFinding.Enqueue(w.Nation[w.IdNations[i]].agents.GetAll()...)
	}
}

type From = common.Vec[int32]
type To = common.Vec[int32]

var cachedPath map[From]map[To][]common.Vec[int32] = map[From]map[To][]common.Vec[int32]{}

func (w *World) setNewPathFindingForWorker(person *Agent) {
}
func (w *World) setNewPathFinding(people ...*Agent) {
	var person *Agent
	var cellTypeToGo CellType
	var posCellsCouldGo []common.Vec[int32]
	var path []common.Vec[int32]
	var cell *BaseCell
	var goal common.Vec[int32]
	var end bool
	//GATHER POSSIBLE PATHS, SORTED BY DISTANCE
	for _, person = range people {
		if person.Status == DEAD {
			continue
		}
		if person.Job == ZOMBIE {
			continue
		}
		cellTypeToGo = JobToCells[person.Job]
		posCellsCouldGo = w.CellType_ToPosCell[cellTypeToGo]
		if len(posCellsCouldGo) == 0 {
			continue
		}
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
				person.Status = IDLE
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
				path = PerformPathFinding_A(w.Map, person.pos, goal)
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
	w.toRunPathFindingForAll()
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
	ps, end := w.personNeedingPathFinding.DenqueueN(10)
	if ps == nil {
		return
	}
	w.setNewPathFinding(ps...)
	if end {
		w.personNeedingPathFinding.Reset()
	}
}

func (w *World) NewNation(idNation ID_NATION, resource map[Resource]float32) {
	w.Nation[idNation] = new(NewNation(w, idNation))
	w.IdNations = append(w.IdNations, idNation)
	w.resources[idNation] = resource
}

func (w *World) NewZombie(where common.Vec[int32]) *Agent {
	return w.Zombies.newAgent(ZOMBIE, where)
}

func (w *World) NewPerson(job Job, where common.Vec[int32], idNation ID_NATION) *Agent {
	if !slices.Contains(w.IdNations, idNation) {
		w.NewNation(idNation, map[Resource]float32{FOOD: 1000})
	}
	p := w.Nation[idNation].newAgent(job, where)
	p.TouchMOVE()
	w.toRunPathFinding(p)
	return p
}

func (w *World) MovementSimulation() (err error) {
	for i := range w.IdNations {
		err = w.Nation[w.IdNations[i]].moveAgentsToGoals()
		if err != nil {
			return err
		}
	}
	err = w.Zombies.moveAgentsToGoals()
	if err != nil {
		return err
	}
	return nil
}

func (w *World) HarvestingSimulation() error {
	for _, idNation := range w.IdNations {
		for celltype, cellsPos := range w.CellType_ToPosCell {
			for _, pos := range cellsPos {
				n := len(w.Nation[idNation].PosToIdAgent[pos])
				for _, q := range CellTypeToResource[celltype] {
					w.resources[idNation][q.What] += float32(n) * q.Amount
				}
			}
		}

		for _, person := range w.Nation[idNation].agents.GetAll() {
			if person.Status == DEAD {
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
		population := w.Nation[idNation].agents.GetAll()
		r := rand.Intn(len(population))
		var person *Agent
		for {
			if maxTime == 0 {
				continue
			}
			person = population[r]
			//TO OPTIMIZE
			if person.Status == DEAD {
				r = rand.Intn(len(population))
				maxTime--
				break
			}
			if person.paths != nil {
				lastPosCell := person.paths.GetLast()
				lastCell, _ := w.Map.GetCell(*lastPosCell)
				lastCell.VirtualNPopulation--
			}
			person.Status = DEAD
			w.Nation[idNation].PosToIdAgent[person.pos] = slices.DeleteFunc(w.Nation[idNation].PosToIdAgent[person.pos], func(a ID_AGENT) bool { return person.id == a })
			w.toRunPathFindingForAll()
			break
		}
	}
	return nil
}
func (w *World) RefreshZombieVision() error {
	peoplePos := []common.Vec[int32]{}
	for _, nation := range w.Nation {
		for pos, l := range nation.PosToIdAgent {
			if len(l) == 0 {
				continue
			}
			peoplePos = append(peoplePos, pos)
		}
	}
	err := w.Zombies.refreshBfsMap(peoplePos)
	if err != nil {
		return err
	}
	return nil
}

func (w *World) zombieEatAgent(agent *Agent) error {
	err := w.Nation[agent.idNation].removeAgent(agent)
	if err != nil {
		return err
	}
	w.Zombies.AddZombie(agent)
	return nil
}

func (w *World) ZombieEat() error {
	for pos, ids := range w.Zombies.PosToIdAgent {
		if len(ids) == 0 {
			continue
		}
		population := []*Agent{}
		for _, nation := range w.Nation {
			for _, id := range nation.PosToIdAgent[pos] {
				agent, err := nation.GetAgent(id)
				if err != nil {
					return err
				}
				if agent.Status == DEAD {
					continue
				}
				population = append(population, agent)
			}
		}
		if len(population) == 0 {
			return nil
		}
		i := rand.Intn(len(population))
		w.zombieEatAgent(population[i])
	}
	return nil
}
