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
	id         ID_NATION
}

func NewNation(w *World, id ID_NATION) Nation {
	return Nation{
		AgentGroup: NewAgentGroup[*Pawn](w),
		id:         id,
		resources:  map[Resource]float32{},
	}
}

func (n *Nation) movePerson(person *Pawn) error {
	from := person.pos
	to, err := person.Move()
	if err != nil {
		return err
	}
	if to == nil {
		return nil
	}
	n.PosToAgents[from] = slices.DeleteFunc(n.PosToAgents[from], func(a *Pawn) bool { return person.id == a.id })
	n.PosToAgents[*to] = append(n.PosToAgents[*to], person)
	return nil
}

func (w *Nation) movePeople() (err error) {
	var person *Pawn
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
		_, err = character.Move()
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
		if person.status == DEAD {
			continue
		}
		for _, q := range JobToConsumingCost[person.job] {
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
		if person.status == DEAD {
			r = rand.Intn(len(population))
			maxTime--
			break
		}
		if person.path != nil {
			lastPosCell := person.path.GetLast()
			lastCell, _ := w.world.CellMap.GetCell(*lastPosCell)
			lastCell.VirtualNPopulation--
		}
		person.status = DEAD
		w.PosToAgents[person.pos] = slices.DeleteFunc(w.PosToAgents[person.pos], func(a *Pawn) bool { return person.id == a.id })
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
