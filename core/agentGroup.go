package core

import (
	"GameOfLife/common"
	"errors"
	"slices"
)

type ID_AGENT int

var CURRENT_ID_AGENT ID_AGENT = 1

type AgentGroup[t Agent] struct {
	world            *World
	agents           *common.SortSlice[t]
	toRunPathFinding bool
	PosToAgents      map[common.Vec[int32]][]t
	id               ID_NATION
}

func (g *AgentGroup[t]) GetAgentsAt(pos common.Vec[int32], condition func(agent t) bool) (res []t) {
	if condition == nil {
		condition = func(_ t) bool { return true }
	}
	for _, agent := range g.PosToAgents[pos] {
		if condition(agent) {
			res = append(res, agent)
		}
	}
	return res
}

func (g *AgentGroup[t]) insertAgent(agent t, where common.Vec[int32]) {
	g.agents.Insert(agent)
	g.PosToAgents[where] = append(g.PosToAgents[where], agent)
}

func (g *AgentGroup[t]) removeAgent(agent t) (err error) {
	removed := g.agents.Remove(agent)
	if !removed {
		return errors.New("Error removing agent")
	}
	g.PosToAgents[agent.Pos()] = slices.DeleteFunc(g.PosToAgents[agent.Pos()], func(a t) bool { return a.Id() == agent.Id() })
	return nil
}
