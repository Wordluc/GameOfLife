package core

import (
	"GameOfLife/common"
	"errors"
	"slices"
)

type ID_AGENT int

var CURRENT_ID_AGENT ID_AGENT = 1

type AgentGroup struct {
	world            *World
	agents           *common.SortSlice[*Agent]
	toRunPathFinding bool
	PosToAgents      map[common.Vec[int32]][]*Agent
	id               ID_NATION
}

func (w *AgentGroup) GetAgentsAt(pos common.Vec[int32], condition func(*Agent) bool) (res []*Agent) {
	if condition == nil {
		condition = func(_ *Agent) bool { return true }
	}

	for i, agent := range w.PosToAgents[pos] {
		if condition(agent) {
			res = append(res, w.PosToAgents[pos][i])
		}
	}
	return res
}

func (w *AgentGroup) newAgent(job Job, where common.Vec[int32]) *Agent {
	p := new(newAgent(job, w.id, where))
	w.agents.Insert(p)
	w.PosToAgents[where] = append(w.PosToAgents[where], p)
	return p
}

func (w *AgentGroup) addAgent(agent *Agent, where common.Vec[int32]) {
	w.agents.Insert(agent)
	agent.pos = where
	if agent.paths != nil {
		cell, _ := w.world.CellMap.GetCell(*agent.paths.GetLast())
		cell.VirtualNPopulation--
	}
	w.PosToAgents[where] = append(w.PosToAgents[where], agent)
}

func (w *Nation) removeAgent(agent *Agent) (err error) {
	removed := w.agents.Remove(agent)
	if !removed {
		return errors.New("Error removing agent")
	}
	w.PosToAgents[agent.pos] = slices.DeleteFunc(w.PosToAgents[agent.pos], func(a *Agent) bool { return a.Id == agent.Id })
	return nil
}
