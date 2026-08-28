package core

import (
	"GameOfLife/common"
	"slices"
)

type ZombieHorde struct {
	AgentGroup
	BfsMap Map[int16]
}

func NewZombieHorde(w *World) ZombieHorde {
	var horde ZombieHorde
	horde = ZombieHorde{
		AgentGroup: AgentGroup{
			world:       w,
			PosToAgents: map[common.Vec[int32]][]*Agent{},
			agents:      common.NewSortSlice(func(a, b *Agent) int { return int(a.id) - int(b.id) }),
		},
		BfsMap: NewMap[int16](w.CellMap.size),
	}
	return horde
}

func (horde *ZombieHorde) moveZombies() (err error) {
	var person *Agent
	for _, person = range slices.Clone(horde.agents.GetAll()) {
		err = horde.moveZombie(person)
		if err != nil {
			return err
		}
	}
	return nil
}

func (horde *ZombieHorde) moveZombie(person *Agent) error {
	if person.IsTouchMOVE() {
		return nil
	}
	neighborhood, _ := horde.BfsMap.GetNeighborhoodCells(person.pos, common.Vec[int32]{X: 3, Y: 3})
	if neighborhood == nil {
		return nil
	}

	cost := neighborhood[person.pos]
	delete(neighborhood, person.pos)
	for pos := range neighborhood {
		if *cost < 0 {
			delete(neighborhood, pos)
		}
		if *neighborhood[pos] < *cost {
			horde.PosToAgents[person.pos] = slices.DeleteFunc(horde.PosToAgents[person.pos], func(a *Agent) bool { return person.id == a.id })
			horde.PosToAgents[pos] = append(horde.PosToAgents[pos], person)
			person.pos = pos
			person.TouchMOVE()
			return nil
		}
	}

	return nil
}

func (z *ZombieHorde) refreshBfsMap(starts []common.Vec[int32]) error {
	m, err := PerformPathFinding_BFS(z.world.CellMap, starts, func(pos common.Vec[int32]) bool {
		if c, err := z.world.CellMap.GetCell(pos); err == nil {
			if slices.Contains([]CellType{WATER, STONE}, c.cellType) {
				return false
			}
		}
		return true
	})
	if err != nil {
		return err
	}
	z.BfsMap = m
	return nil
}

func (z *ZombieHorde) addZombie(agent *Agent) error {
	if agent.paths != nil {
		cell, _ := z.world.CellMap.GetCell(*agent.paths.GetLast())
		cell.VirtualNPopulation--
	}
	agent.Status = MOVING
	agent.paths = nil
	agent.idNation = -1
	z.addAgent(agent, agent.pos)
	return nil
}
