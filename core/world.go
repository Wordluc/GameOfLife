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
	CellMap            *Map[BaseCell]
	AgentsMap          *Map[[]*Agent]
	Nations            map[ID_NATION]*Nation
	Zombies            *ZombieHorde
	CellType_ToPosCell bindingCellTypeToCell
	IdNations          []ID_NATION

	personNeedingPathFinding *common.Queue[*Agent]
}

func NewWorld(size common.Vec[int32]) (w *World) {
	w = &World{}
	w.CellMap = new(NewMap[BaseCell](size))
	w.AgentsMap = new(NewMap[[]*Agent](size))
	w.Nations = map[ID_NATION]*Nation{}
	w.CellType_ToPosCell = make(bindingCellTypeToCell)
	w.personNeedingPathFinding = common.NewQueue[*Agent](nil, nil)
	w.Zombies = new(NewZombieHorde(w))
	return w
}

func (w *World) GenerateMap() {
	var c *BaseCell
	w.CellMap.SetEachElement(func(x, y int32) *BaseCell {
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
		w.personNeedingPathFinding.Enqueue(w.Nations[w.IdNations[i]].agents.GetAll()...)
	}
}

type From = common.Vec[int32]
type To = common.Vec[int32]

var cachedPath map[From]map[To][]common.Vec[int32] = map[From]map[To][]common.Vec[int32]{}

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
				cell, _ = w.CellMap.GetCell(*person.paths.GetLast())
				cell.VirtualNPopulation--
				person.paths = nil
			}

			if cachedPath[person.pos][goal] != nil {
				path = cachedPath[person.pos][goal]
			} else {
				path = PerformPathFinding_A(w.CellMap, person.pos, goal, func(pos common.Vec[int32]) bool {
					agents, _ := w.AgentsMap.GetCell(pos)
					if agents != nil && len(*agents) != 0 && slices.ContainsFunc(*agents, func(a *Agent) bool { return a.Job != person.Job && a.Status == WORKING }) {
						return false
					}
					if c, err := w.CellMap.GetCell(pos); err == nil {
						if slices.Contains([]CellType{WATER, STONE}, c.cellType) {
							return false
						}
					}
					return true
				})
				cachedPath[person.pos][goal] = path
			}

			if len(path) != 0 {
				cell, _ = w.CellMap.GetCell(goal)
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
	if pos.X-halfX < 0 || pos.Y-halfY < 0 || pos.X+halfX >= w.CellMap.size.X || pos.Y+halfY >= w.CellMap.size.Y {
		return errors.New("Out of bound")
	}

	neighborhood, err := w.CellMap.GetNeighborhoodCells(pos, size)
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
			if !definition.CanConvert(pos, w.CellMap) {
				return fmt.Errorf("CanConvert %v failed", cellType)
			}
		}
	}
	for pos := range neighborhood {
		cell, err := w.CellMap.GetCell(pos)
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
		cell, err = w.CellMap.GetCell(pos)
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

func (w *World) AddNation(idNation ID_NATION, resource map[Resource]float32) {
	n := NewNation(w, idNation)
	n.resources = resource
	w.Nations[idNation] = new(n)
	w.IdNations = append(w.IdNations, idNation)
}

func (w *World) NewZombie(where common.Vec[int32]) *Agent {
	return w.Zombies.newAgent(ZOMBIE, where)
}

func (w *World) NewPerson(job Job, where common.Vec[int32], idNation ID_NATION) *Agent {
	if !slices.Contains(w.IdNations, idNation) {
		w.AddNation(idNation, map[Resource]float32{FOOD: 1000})
	}
	p := w.Nations[idNation].newAgent(job, where)
	p.TouchMOVE()
	w.toRunPathFinding(p)
	return p
}

func (w *World) MovementSimulation() (err error) {
	w.AgentsMap = new(NewMap[[]*Agent](w.AgentsMap.size))
	for _, nation := range w.Nations {
		err = nation.movePeople()
		if err != nil {
			return err
		}
		err = nation.moveCharacters()
		if err != nil {
			return err
		}
		for pos, agents := range nation.PosToAgents {
			w.AgentsMap.SetRawCell(new(agents), pos)
		}
	}
	err = w.Zombies.moveZombies()
	if err != nil {
		return err
	}
	for pos, agents := range w.Zombies.PosToAgents {
		w.AgentsMap.SetRawCell(new(agents), pos)
	}
	return nil
}

func (w *World) HarvestingSimulation() error {
	for _, nation := range w.Nations {
		nation.Harvesting()
	}
	return nil
}

func (w *World) StarvingSimulation() error {
	for _, nation := range w.Nations {
		nation.Starving()
	}
	return nil
}

func (w *World) RefreshZombieVision() error {
	zombieGoals := []common.Vec[int32]{}
	for _, nation := range w.Nations {
		for pos, l := range nation.PosToAgents {
			if len(l) == 0 {
				continue
			}
			zombieGoals = append(zombieGoals, pos)
		}
		for _, c := range nation.Characters {
			zombieGoals = append(zombieGoals, c.pos)
		}
	}
	return w.Zombies.refreshBfsMap(zombieGoals)
}

func (w *World) zombieEatAgent(agent *Agent) error {
	err := w.Nations[agent.idNation].removeAgent(agent)
	if err != nil {
		return err
	}
	w.Zombies.addZombie(agent)
	return nil
}

func (w *World) ZombieEat() error {
	for pos, ids := range w.Zombies.PosToAgents {
		if len(ids) == 0 {
			continue
		}
		population := []*Agent{}
		for _, nation := range w.Nations {
			for _, agent := range nation.PosToAgents[pos] {
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

func (w *World) GetCharactersAt(pos common.Vec[int32], idNation *ID_NATION) (res []*Character) {
	if idNation == nil {
		for _, nation := range w.Nations {
			res = append(res, nation.GetCharactersAt(pos)...)
		}
		return res
	}
	if w.Nations == nil {
		return nil
	}
	return w.Nations[*idNation].Characters
}

func (w *World) SetPathCharactersTo(end common.Vec[int32], characters ...*Character) (err error) {
	if characters == nil {
		return errors.New("No character selected")
	}
	for i := range characters {
		if characters[i].pos.IsEqual(end) {
			continue
		}
		characters[i].paths = common.NewQueue(PerformPathFinding_A(w.CellMap, characters[i].pos, end, func(pos common.Vec[int32]) bool {
			if c, err := w.CellMap.GetCell(pos); err == nil {
				if slices.Contains([]CellType{WATER, STONE}, c.cellType) {
					return false
				}
			}
			return true
		}), nil)
		characters[i].paths.Denqueue()
		pos, _ := characters[i].paths.Denqueue()
		characters[i].pos = pos
	}
	return nil
}

func (w *World) GetAgentsAt(pos common.Vec[int32], zombie bool) (res []*Agent) {
	if zombie {
		return w.Zombies.GetAgentsAt(pos, nil)
	}
	for _, nation := range w.Nations {
		res = append(res, nation.GetAgentsAt(pos, nil)...)
	}
	return res
}
