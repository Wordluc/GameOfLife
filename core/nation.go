package core

import (
	"GameOfLife/common"
	"errors"
	"slices"
)

type Nation struct {
	AgentGroup
}

func NewNation(w *World, id ID_NATION) Nation {
	return Nation{
		AgentGroup: AgentGroup{
			World:        w,
			PosToIdAgent: map[common.Vec[int32]][]ID_AGENT{},
			agents:       common.NewSortSlice(func(a, b *Agent) int { return int(a.id) - int(b.id) }),
			id:           id,
		},
	}
}

func (n *Nation) movePersonToGoalUsingPathFindingA(person *Agent) error {
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

func (w *Nation) moveAgentsToGoals() (err error) {
	var person *Agent
	for _, person = range slices.Clone(w.agents.GetAll()) {
		err = w.movePersonToGoalUsingPathFindingA(person)
		if err != nil {
			return err
		}
	}
	return nil
}
