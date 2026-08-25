package core

import (
	"GameOfLife/common"
	"errors"
	"slices"
)

type ID_AGENT int

var CURRENT_ID_AGENT ID_AGENT = 1

type AgentGroup struct {
	*World
	agents           *common.SortSlice[*Agent]
	toRunPathFinding bool
	PosToIdAgent     map[common.Vec[int32]][]ID_AGENT
	id               ID_NATION
}

func (w *AgentGroup) GetAgent(id ID_AGENT) (res *Agent, err error) {
	err = w.agents.ForEach(func(index int, value *Agent) (stop bool, err error) {
		if value.id == id {
			res = value
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return res, err
	}
	return res, nil
}

func (w *AgentGroup) GetAgentsIdInsideCellf(pos common.Vec[int32], check func(a *Agent) bool) (res []ID_AGENT) {
	var agent *Agent
	var err error
	for _, id := range w.PosToIdAgent[pos] {
		agent, err = w.GetAgent(id)
		if err != nil {
			return nil
		}
		if !check(agent) {
			continue
		}
		res = append(res, agent.id)
	}
	return res
}

func (w *AgentGroup) GetAgentsIdInsideCell(pos common.Vec[int32]) (res []ID_AGENT) {
	res = append(res, w.PosToIdAgent[pos]...)
	return res
}

func (w *AgentGroup) newAgent(job Job, where common.Vec[int32]) *Agent {
	p := new(newPerson(job, w.id, where))
	w.agents.Insert(p)
	w.PosToIdAgent[where] = append(w.PosToIdAgent[where], p.id)
	return p
}

func (w *AgentGroup) addAgent(agent *Agent, where common.Vec[int32]) {
	w.agents.Insert(agent)
	agent.pos = where
	if agent.paths != nil {
		cell, _ := w.Map.GetCell(*agent.paths.GetLast())
		cell.VirtualNPopulation--
	}
	w.PosToIdAgent[where] = append(w.PosToIdAgent[where], agent.id)
}

func (w *Nation) removeAgent(agent *Agent) (err error) {
	removed := w.agents.Remove(agent)
	if !removed {
		return errors.New("Error removing agent")
	}
	w.PosToIdAgent[agent.pos] = slices.DeleteFunc(w.PosToIdAgent[agent.pos], func(a ID_AGENT) bool { return a == agent.id })
	return nil
}
