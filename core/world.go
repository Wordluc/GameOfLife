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
	agentsMap          *Map[[]Agent]
	Nations            map[ID_NATION]*Nation
	Zombies            *ZombieHorde
	cellType_ToPosCell bindingCellTypeToCell
	IdNations          []ID_NATION

	personNeedingPathFinding *common.Queue[*Pawn]
}

func NewWorld(size common.Vec[int32]) (w *World) {
	w = &World{}
	w.CellMap = new(NewMap[BaseCell](size))
	w.agentsMap = new(NewMap[[]Agent](size))
	w.Nations = map[ID_NATION]*Nation{}
	w.cellType_ToPosCell = make(bindingCellTypeToCell)
	w.personNeedingPathFinding = common.NewQueue[*Pawn](nil, func(a, b *Pawn) int { return int(a.Id()) - int(b.Id()) })
	w.Zombies = new(NewZombieHorde(w))
	return w
}

func (w *World) GenerateMap() {
	var c *BaseCell
	w.CellMap.SetEachElement(func(x, y int32) *BaseCell {
		c = NewEmptyBaseCell(common.Vec[int32]{X: x, Y: y})
		w.cellType_ToPosCell.SetCellTypeCell(c, GRASS)
		return c
	})
}

func (w *World) toRunPathFinding(ps ...*Pawn) {
	w.personNeedingPathFinding.EnqueueUnique(ps...)
}

func (w *World) toRunPathFindingForAll() {
	for i := range w.IdNations {
		for _, pawn := range w.Nations[w.IdNations[i]].agents.GetAll() {
			w.personNeedingPathFinding.EnqueueUnique(pawn)
		}
	}
}

type From = common.Vec[int32]
type To = common.Vec[int32]

func (w *World) setNewPathFinding(people ...*Pawn) {
	var cachedPath map[From]map[To][]common.Vec[int32] = map[From]map[To][]common.Vec[int32]{}
	var person *Pawn
	var cellTypeToGo CellType
	var posCellsCouldGo []common.Vec[int32]
	var path []common.Vec[int32]
	var cell *BaseCell
	var goal common.Vec[int32]
	var end bool
	//GATHER POSSIBLE PATHS, SORTED BY DISTANCE
	for _, person = range people {
		if person.Status() == DEAD {
			continue
		}
		cellTypeToGo = JobToCells[person.Job()]
		posCellsCouldGo = w.cellType_ToPosCell[cellTypeToGo]
		if len(posCellsCouldGo) == 0 {
			continue
		}
		goals := common.NewQueue(
			posCellsCouldGo,
			func(a, b common.Vec[int32]) int {
				return cmp.Compare(fromAtoB(person.Pos(), a), fromAtoB(person.Pos(), b))
			})
		if cachedPath[person.Pos()] == nil {
			cachedPath[person.Pos()] = make(map[To][]common.Vec[int32])
		}
		for {
			goal, end = goals.Denqueue()
			if end {
				person.SetStatus(IDLE)
				w.personNeedingPathFinding.Enqueue(person)
				break
			}
			if goal.IsEqual(person.Pos()) {
				break
			}

			if person.Paths() != nil {
				cell, _ = w.CellMap.GetCell(*person.Paths().GetLast())
				cell.VirtualNPopulation--
				person.SetPaths(nil)
			}

			if cachedPath[person.Pos()][goal] != nil {
				path = cachedPath[person.Pos()][goal]
			} else {
				path = PerformPathFinding_A(w.CellMap, person.Pos(), goal, func(pos common.Vec[int32]) bool {
					agents, _ := w.agentsMap.GetCell(pos)
					if agents != nil && len(*agents) != 0 && slices.ContainsFunc(*agents, func(a Agent) bool {
						if pawn, ok := a.(*Pawn); ok {
							return pawn.Job() != person.Job()
						}
						return true
					}) {
						return false
					}
					if c, err := w.CellMap.GetCell(pos); err == nil {
						if slices.Contains([]CellType{WATER, STONE}, c.cellType) {
							return false
						}
					}
					return true
				})
				cachedPath[person.Pos()][goal] = path
			}

			if len(path) != 0 {
				cell, _ = w.CellMap.GetCell(goal)
				if cell.VirtualNPopulation >= cellsDefinition[cell.cellType].maxPeople {
					continue
				}
				cell.VirtualNPopulation++
				person.SetPaths(common.NewQueue(path, nil))
				//REMOVE FIRST ELEMENT, THE ORIGIN (person.pos)
				person.Paths().Denqueue()
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
		w.cellType_ToPosCell.SetCellTypeCell(cell, cellType)
		err = definition.convert(cell)
		if err != nil {
			return err
		}
	}
	w.toRunPathFindingForAll()
	return nil
}

func (w *World) GetCellsByType(cellType CellType) (res []*BaseCell, err error) {
	var cell *BaseCell
	for _, pos := range w.cellType_ToPosCell[cellType] {
		cell, err = w.CellMap.GetCell(pos)
		if err != nil {
			return nil, err
		}
		res = append(res, cell)
	}
	return res, nil
}

func (w *World) PerformPathFinding() {
	ps, _ := w.personNeedingPathFinding.DenqueueN(10)
	if ps == nil {
		return
	}
	w.setNewPathFinding(ps...)
	if w.personNeedingPathFinding.Remaining() <= 0 {
		w.personNeedingPathFinding.Reset()
	}
}

func (w *World) AddNation(idNation ID_NATION, resource map[Resource]float32) {
	n := NewNation(w, idNation)
	n.resources = resource
	w.Nations[idNation] = new(n)
	w.IdNations = append(w.IdNations, idNation)
}

func (w *World) NewZombie(where common.Vec[int32]) *Zombie {
	return w.Zombies.newZombie(where)
}

func (w *World) NewPerson(job Job, where common.Vec[int32], idNation ID_NATION) *Pawn {
	if !slices.Contains(w.IdNations, idNation) {
		w.AddNation(idNation, map[Resource]float32{FOOD: 1000})
	}
	p := w.Nations[idNation].newPawn(job, where)
	w.addToAgentsMap(p, nil)
	w.toRunPathFinding(p)
	return p
}

func (w *World) MovementSimulation() (err error) {
	w.toRunPathFindingForAll()
	for _, nation := range w.Nations {
		err = nation.movePeople()
		if err != nil {
			return err
		}
		err = nation.moveCharacters()
		if err != nil {
			return err
		}
	}
	err = w.Zombies.moveZombies()
	if err != nil {
		return err
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
			zombieGoals = append(zombieGoals, c.Pos())
		}
	}
	return w.Zombies.refreshBfsMap(zombieGoals)
}

func (w *World) removeFromAgentsMap(agent Agent, pos *common.Vec[int32]) {
	if pos == nil {
		pos = new(agent.Pos())
	}
	agents, _ := w.agentsMap.GetCell(*pos)
	if agents == nil {
		return
	}
	*agents = slices.DeleteFunc(*agents, func(a Agent) bool { return a.Id() == agent.Id() })
	w.agentsMap.SetRawCell(agents, *pos)
}

func (w *World) addToAgentsMap(agent Agent, pos *common.Vec[int32]) {
	if pos == nil {
		pos = new(agent.Pos())
	}
	agents, _ := w.agentsMap.GetCell(*pos)
	if agents == nil {
		agents = new([]Agent)
	}
	*agents = append(*agents, agent)
	w.agentsMap.SetRawCell(agents, *pos)
}

func (w *World) zombieEatAgent(agent *Pawn, spawnAt common.Vec[int32]) error {
	err := w.Nations[agent.IdNation()].removeAgent(agent)
	if err != nil {
		return err
	}
	w.removeFromAgentsMap(agent, nil)
	zombie, err := w.Zombies.addZombie(agent, spawnAt)
	if err != nil {
		return err
	}
	w.addToAgentsMap(zombie, nil)
	return nil
}

func (w *World) ZombieEat() error {
	for pos, ids := range w.Zombies.PosToAgents {
		if len(ids) == 0 {
			continue
		}
		population := []Agent{}
		//TODO: define for each agent a attack range
		n, _ := w.agentsMap.GetNeighborhoodCells(pos, common.Vec[int32]{X: 3, Y: 3})
		if n == nil {
			continue
		}
		for _, agents := range n {
			if agents == nil {
				continue
			}
			for _, agent := range *agents {
				if agent == nil {
					continue
				}
				if agent.Status() == DEAD {
					continue
				}
				if isZombie(agent) {
					continue
				}
				if isCharacter(agent) {
					continue
				}
				population = append(population, agent)
			}
		}
		if len(population) == 0 {
			return nil
		}
		i := rand.Intn(len(population))

		pawn, ok := population[i].(*Pawn)
		if !ok {
			return nil
		}
		w.zombieEatAgent(pawn, pos)
		println("loop")
	}
	return nil
}

func (w *World) AddCharactersAt(pos common.Vec[int32], idNation ID_NATION) error {
	if int(idNation) >= len(w.Nations) {
		return errors.New("Missing nation")
	}
	character := NewCharacter(pos)
	w.Nations[idNation].Characters = append(w.Nations[idNation].Characters, character)
	w.addToAgentsMap(character, nil)
	w.toRunPathFindingForAll()
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
		current := characters[i].Pos()
		if current.IsEqual(end) {
			continue
		}
		characters[i].SetPaths(common.NewQueue(PerformPathFinding_A(w.CellMap, current, end, func(pos common.Vec[int32]) bool {
			if c, err := w.CellMap.GetCell(pos); err == nil {
				if slices.Contains([]CellType{WATER, STONE}, c.cellType) {
					return false
				}
			}
			agents, _ := w.agentsMap.GetCell(pos)
			if agents != nil && len(*agents) != 0 && slices.ContainsFunc(*agents, func(a Agent) bool { return !isCharacter(a) }) {
				return false
			}

			return true
		}), nil))
		characters[i].Paths().Denqueue()
	}
	return nil
}

func (w *World) GetAgentsAt(pos common.Vec[int32], idNation *ID_NATION, zombie bool) (res []Agent) {
	if zombie {
		for _, z := range w.Zombies.GetAgentsAt(pos, nil) {
			res = append(res, z)
		}
		return res
	}
	if idNation == nil {
		for _, nation := range w.Nations {
			for _, pawn := range nation.GetAgentsAt(pos, nil) {
				res = append(res, pawn)
			}
		}
		return res
	}
	for _, pawn := range w.Nations[*idNation].GetAgentsAt(pos, nil) {
		res = append(res, pawn)
	}
	return res
}
