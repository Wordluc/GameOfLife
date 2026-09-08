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

func (n *Nation) movePerson(person *Agent) error {
	from := person.pos
	to, err := person.FollowPath_StarA()
	if err != nil {
		return nil
	}
	if to == nil {
		return nil
	}
	n.PosToAgents[from] = slices.DeleteFunc(n.PosToAgents[from], func(a *Agent) bool { return person.Id == a.Id })
	n.PosToAgents[*to] = append(n.PosToAgents[*to], person)
	person.pos = *to
	return nil
}

func (w *Nation) movePeople() (err error) {
	var person *Agent
	for _, person = range w.agents.GetAll() {
		err = w.movePerson(person)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *Nation) moveCharacters() (err error) {
	var character *Character
	for _, character = range w.Characters {
		_, err = character.FollowPath_StarA()
		if err != nil {
			return err
		}
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
