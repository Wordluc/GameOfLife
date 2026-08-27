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
			world:        w,
			PosToIdAgent: map[common.Vec[int32]][]ID_AGENT{},
			agents:       common.NewSortSlice(func(a, b *Agent) int { return int(a.id) - int(b.id) }),
		},
		BfsMap: NewMap[int16](w.Map.size),
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
	for i := range neighborhood {
		if *cost < 0 {
			delete(neighborhood, i)
		}
		if *neighborhood[i] < *cost {
			horde.PosToIdAgent[person.pos] = slices.DeleteFunc(horde.PosToIdAgent[person.pos], func(a ID_AGENT) bool { return person.id == a })
			horde.PosToIdAgent[i] = append(horde.PosToIdAgent[i], person.id)
			person.pos = i
			person.TouchMOVE()
			return nil
		}
	}

	return nil
}

func (z *ZombieHorde) refreshBfsMap(starts []common.Vec[int32]) error {
	m, err := PerformPathFinding_BFS(z.world.Map, starts)
	if err != nil {
		return err
	}
	z.BfsMap = m
	return nil
}

func (z *ZombieHorde) addZombie(agent *Agent) error {
	if agent.paths != nil {
		cell, _ := z.world.Map.GetCell(*agent.paths.GetLast())
		cell.VirtualNPopulation--
	}
	agent.Status = MOVING
	agent.paths = nil
	agent.idNation = -1
	z.addAgent(agent, agent.pos)
	return nil
}
