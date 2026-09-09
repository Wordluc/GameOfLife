package core

import (
	"GameOfLife/common"
	"math/rand"
	"slices"
)

type Nation struct {
	AgentGroup
	resources  map[Resource]float32
	Characters []*Character
}

func NewNation(w *World, id ID_NATION) Nation {
	return Nation{
		AgentGroup: AgentGroup{
			world:       w,
			PosToAgents: map[common.Vec[int32]][]*Agent{},
			agents:      common.NewSortSlice(func(a, b *Agent) int { return int(a.Id) - int(b.Id) }),
			id:          id,
		},
		resources: map[Resource]float32{},
	}
}

func (n *Nation) movePerson(person *Agent) (from common.Vec[int32], to *common.Vec[int32], err error) {
	from = person.pos
	to, err = person.FollowPath_StarA()
	if err != nil {
		return from, to, err
	}
	if to == nil {
		return from, to, err
	}
	n.PosToAgents[from] = slices.DeleteFunc(n.PosToAgents[from], func(a *Agent) bool { return person.Id == a.Id })
	n.PosToAgents[*to] = append(n.PosToAgents[*to], person)
	person.pos = *to
	return from, to, err
}

func (w *Nation) movePeople(agentMap *Map[[]*Agent]) (err error) {
	var person *Agent
	for _, person = range w.agents.GetAll() {
		from, to, err := w.movePerson(person)
		if err != nil {
			return err
		}
		if to == nil {
			continue
		}
		agents, _ := agentMap.GetCell(from)
		if agents == nil {
			agents = new([]*Agent)
		}
		*agents = slices.DeleteFunc(*agents, func(a *Agent) bool { return a.Id == person.Id })
		agentMap.SetRawCell(agents, from)

		agents, _ = agentMap.GetCell(*to)
		if agents == nil {
			agents = new([]*Agent)
		}
		agentMap.SetRawCell(new(append(*agents, person)), *to)

	}
	return nil
}

func (n *Nation) moveCharacter(character *Character) (from common.Vec[int32], to *common.Vec[int32], err error) {
	from = character.pos
	to, err = character.FollowPath_StarA()
	return from, to, err
}

func (w *Nation) moveCharacters(agentMap *Map[[]*Agent]) (err error) {
	var character *Character
	for _, character = range w.Characters {
		from, to, err := w.moveCharacter(character)
		if err != nil {
			return err
		}
		if to == nil {
			continue
		}
		agents, _ := agentMap.GetCell(from)
		if agents == nil {
			agents = new([]*Agent)
		}
		*agents = slices.DeleteFunc(*agents, func(a *Agent) bool { return a.Id == character.Id })
		agentMap.SetRawCell(agents, from)

		agents, _ = agentMap.GetCell(*to)
		if agents == nil {
			agents = new([]*Agent)
		}
		*agents = append(*agents, &character.Agent)
		agentMap.SetRawCell(agents, *to)
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
		if person.Status == DEAD {
			continue
		}
		for _, q := range JobToConsumingCost[person.Job] {
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
			lastCell, _ := w.world.CellMap.GetCell(*lastPosCell)
			lastCell.VirtualNPopulation--
		}
		person.Status = DEAD
		w.PosToAgents[person.pos] = slices.DeleteFunc(w.PosToAgents[person.pos], func(a *Agent) bool { return person.Id == a.Id })
		w.world.toRunPathFindingForAll()
		break
	}
	return nil
}

func (w *Nation) GetCharactersAt(pos common.Vec[int32]) (res []*Character) {
	for i := range w.Characters {
		if w.Characters[i].pos.IsEqual(pos) {
			res = append(res, w.Characters[i])
		}
	}
	return res
}
