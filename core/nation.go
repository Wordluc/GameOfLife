package core

import (
	"GameOfLife/common"
	"errors"
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
			world:        w,
			PosToIdAgent: map[common.Vec[int32]][]ID_AGENT{},
			agents:       common.NewSortSlice(func(a, b *Agent) int { return int(a.id) - int(b.id) }),
			id:           id,
		},
		resources: map[Resource]float32{},
	}
}

func (n *Nation) movePerson(person *Agent) error {
	if person.IsTouchMOVE() {
		return nil
	}
	if person.paths == nil {
		return nil
	}
	from := person.paths.GetBack(1)
	if from != nil && !from.IsEqual(person.pos) {
		return errors.New("Error initial position")
	}

	to, end := person.paths.Denqueue()
	if end {
		person.Status = WORKING
		return nil
	}
	person.Status = MOVING
	n.PosToIdAgent[*from] = slices.DeleteFunc(n.PosToIdAgent[*from], func(a ID_AGENT) bool { return person.id == a })
	n.PosToIdAgent[to] = append(n.PosToIdAgent[to], person.id)
	person.pos = to
	person.TouchMOVE()
	return nil
}

func (n *Nation) moveCharacter(playable *Character) error {
	if playable.paths == nil {
		return nil
	}
	from := playable.paths.GetBack(1)
	if from != nil && !from.IsEqual(playable.pos) {
		return errors.New("Error initial position")
	}

	to, end := playable.paths.Denqueue()
	if end {
		return nil
	}
	playable.pos = to
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
		err = w.moveCharacter(character)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *Nation) Harvesting() error {
	for celltype, cellsPos := range w.world.CellType_ToPosCell {
		for _, pos := range cellsPos {
			n := len(w.PosToIdAgent[pos])
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
			lastCell, _ := w.world.Map.GetCell(*lastPosCell)
			lastCell.VirtualNPopulation--
		}
		person.Status = DEAD
		w.PosToIdAgent[person.pos] = slices.DeleteFunc(w.PosToIdAgent[person.pos], func(a ID_AGENT) bool { return person.id == a })
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
