package core

import (
	"GameOfLife/common"
	"slices"
)

type Zombie struct {
	agentCore
}

func newZombie(where common.Vec[int32]) *Zombie {
	return &Zombie{
		agentCore: newAgentCore(where),
	}
}

func isZombie(a Agent) bool {
	_, ok := a.(*Zombie)
	return ok
}

type ZombieHorde struct {
	AgentGroup[*Zombie]
	BfsMap Map[int16]
}

func NewZombieHorde(w *World) ZombieHorde {
	var horde ZombieHorde
	horde = ZombieHorde{
		AgentGroup: AgentGroup[*Zombie]{
			world:       w,
			PosToAgents: map[common.Vec[int32]][]*Zombie{},
			agents:      common.NewSortSlice(func(a, b *Zombie) int { return int(a.Id()) - int(b.Id()) }),
		},
		BfsMap: NewMap[int16](w.CellMap.size),
	}
	return horde
}

func (horde *ZombieHorde) newZombie(where common.Vec[int32]) *Zombie {
	zombie := newZombie(where)
	horde.insertAgent(zombie, where)
	return zombie
}

func (horde *ZombieHorde) addZombie(pawn *Pawn, spawnAt common.Vec[int32]) (*Zombie, error) {
	if pawn.Paths() != nil {
		cell, _ := horde.world.CellMap.GetCell(*pawn.Paths().GetLast())
		cell.VirtualNPopulation--
	}
	pawn.SetStatus(MOVING)
	pawn.SetPaths(nil)
	pawn.SetPos(spawnAt)
	zombie := &Zombie{agentCore: pawn.agentCore}
	horde.insertAgent(zombie, spawnAt)
	return zombie, nil
}

func (horde *ZombieHorde) moveZombies() (err error) {
	var zombie *Zombie
	for _, zombie = range slices.Clone(horde.agents.GetAll()) {
		from, to, err := horde.moveZombie(zombie)
		if err != nil {
			return err
		}
		if to == nil {
			continue
		}
		horde.world.removeFromAgentsMap(zombie, &from)
		horde.world.addToAgentsMap(zombie, to)
	}
	return nil
}

func (horde *ZombieHorde) moveZombie(zombie *Zombie) (from common.Vec[int32], to *common.Vec[int32], err error) {
	from = zombie.Pos()
	neighborhood, _ := horde.BfsMap.GetNeighborhoodCells(zombie.Pos(), common.Vec[int32]{X: 3, Y: 3})
	if neighborhood == nil {
		return from, to, nil
	}

	cost := neighborhood[zombie.Pos()]
	delete(neighborhood, zombie.Pos())
	for pos := range neighborhood {
		if *cost < 0 {
			delete(neighborhood, pos)
		}
		if *neighborhood[pos] < *cost {
			if len(horde.world.GetAgentsAt(pos, nil, false)) != 0 {
				continue
			}
			if len(horde.world.GetCharactersAt(pos, nil)) != 0 {
				continue
			}
			horde.PosToAgents[zombie.Pos()] = slices.DeleteFunc(horde.PosToAgents[zombie.Pos()], func(a *Zombie) bool { return zombie.Id() == a.Id() })
			horde.PosToAgents[pos] = append(horde.PosToAgents[pos], zombie)
			zombie.SetPos(pos)
			return from, &zombie.agentCore.pos, nil
		}
	}
	//RANDOM MOVEMENT
	for key := range neighborhood {
		if *neighborhood[key] == *cost && !key.IsEqual(zombie.Pos()) {
			horde.PosToAgents[zombie.Pos()] = slices.DeleteFunc(horde.PosToAgents[zombie.Pos()], func(a *Zombie) bool { return zombie.Id() == a.Id() })
			horde.PosToAgents[key] = append(horde.PosToAgents[key], zombie)
			zombie.SetPos(key)
		}
	}

	return from, to, nil
}

func (z *ZombieHorde) refreshBfsMap(starts []common.Vec[int32]) error {
	m, err := PerformPathFinding_BFS(z.world.CellMap, starts, func(pos common.Vec[int32]) bool {
		if c, err := z.world.CellMap.GetCell(pos); err == nil {
			if slices.Contains([]CellType{WATER, STONE}, c.cellType) {
				return false
			}
			agents, _ := z.world.agentsMap.GetCell(pos)
			if agents != nil && len(*agents) != 0 && slices.ContainsFunc(*agents, func(a Agent) bool { return !isZombie(a) }) {
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
