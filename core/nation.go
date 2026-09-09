package core

import (
	"GameOfLife/common"
	"math/rand"
	"slices"
)

type Nation struct {
	AgentGroup[*Pawn]
	resources  map[Resource]float32
	Characters []*Character
}

func NewNation(w *World, id ID_NATION) Nation {
	return Nation{
		AgentGroup: AgentGroup[*Pawn]{
			world:       w,
			PosToAgents: map[common.Vec[int32]][]*Pawn{},
			agents:      common.NewSortSlice(func(a, b *Pawn) int { return int(a.Id()) - int(b.Id()) }),
			id:          id,
		},
		resources: map[Resource]float32{},
	}
}

func (n *Nation) newPawn(job Job, where common.Vec[int32]) *Pawn {
	pawn := NewPawn(job, n.id, where)
	n.insertAgent(pawn, where)
	return pawn
}

func (n *Nation) movePerson(person *Pawn) (from common.Vec[int32], to *common.Vec[int32], err error) {
	from = person.Pos()
	to, err = person.FollowPath_StarA()
	if err != nil {
		return from, to, err
	}
	if to == nil {
		return from, to, err
	}
	n.PosToAgents[from] = slices.DeleteFunc(n.PosToAgents[from], func(a *Pawn) bool { return person.Id() == a.Id() })
	n.PosToAgents[*to] = append(n.PosToAgents[*to], person)
	person.SetPos(*to)
	return from, to, err
}

func (w *Nation) movePeople() (err error) {
	var person *Pawn
	for _, person = range w.agents.GetAll() {
		from, to, err := w.movePerson(person)
		if err != nil {
			return err
		}
		if to == nil {
			continue
		}
		agents, _ := w.world.agentsMap.GetCell(from)
		if agents == nil {
			agents = new([]Agent)
		}
		*agents = slices.DeleteFunc(*agents, func(a Agent) bool { return a.Id() == person.Id() })
		w.world.agentsMap.SetRawCell(agents, from)

		agents, _ = w.world.agentsMap.GetCell(*to)
		if agents == nil {
			agents = new([]Agent)
		}
		*agents = append(*agents, person)
		w.world.agentsMap.SetRawCell(agents, *to)

	}
	return nil
}

func (n *Nation) moveCharacter(character *Character) (from common.Vec[int32], to *common.Vec[int32], err error) {
	from = character.Pos()
	to, err = character.FollowPath_StarA()
	return from, to, err
}

func (w *Nation) moveCharacters() (err error) {
	var character *Character
	for _, character = range w.Characters {
		from, to, err := w.moveCharacter(character)
		if err != nil {
			return err
		}
		if to == nil {
			continue
		}
		agents, _ := w.world.agentsMap.GetCell(from)
		if agents == nil {
			agents = new([]Agent)
		}
		*agents = slices.DeleteFunc(*agents, func(a Agent) bool { return a.Id() == character.Id() })
		w.world.agentsMap.SetRawCell(agents, from)

		agents, _ = w.world.agentsMap.GetCell(*to)
		if agents == nil {
			agents = new([]Agent)
		}
		*agents = append(*agents, character)
		w.world.agentsMap.SetRawCell(agents, *to)
	}
	return nil
}

func (w *Nation) Harvesting() error {
	for celltype, cellsPos := range w.world.cellType_ToPosCell {
		for _, pos := range cellsPos {
			n := len(w.PosToAgents[pos])
			for _, q := range CellTypeToResource[celltype] {
				w.resources[q.What] += float32(n) * q.Amount
			}
		}
	}

	for _, person := range w.agents.GetAll() {
		if person.Status() == DEAD {
			continue
		}
		for _, q := range JobToConsumingCost[person.Job()] {
			w.resources[q.What] -= q.Amount
		}
	}
	return nil
}

func (w *Nation) Starving() error {
	if w.resources[FOOD] > -100 {
		return nil
	}
	var maxTime = 10
	population := w.agents.GetAll()
	r := rand.Intn(len(population))
	var person *Pawn
	for {
		if maxTime == 0 {
			continue
		}
		person = population[r]
		//TO OPTIMIZE
		if person.Status() == DEAD {
			r = rand.Intn(len(population))
			maxTime--
			break
		}
		if person.Paths() != nil {
			lastPosCell := person.Paths().GetLast()
			lastCell, _ := w.world.CellMap.GetCell(*lastPosCell)
			lastCell.VirtualNPopulation--
		}
		person.SetStatus(DEAD)
		w.PosToAgents[person.Pos()] = slices.DeleteFunc(w.PosToAgents[person.Pos()], func(a *Pawn) bool { return person.Id() == a.Id() })
		w.world.toRunPathFindingForAll()
		break
	}
	return nil
}

func (w *Nation) GetCharactersAt(pos common.Vec[int32]) (res []*Character) {
	for i := range w.Characters {
		p := w.Characters[i].Pos()
		if p.IsEqual(pos) {
			res = append(res, w.Characters[i])
		}
	}
	return res
}
