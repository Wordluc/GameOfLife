package core

import (
	"GameOfLife/common"
	"errors"
	"slices"
)

var CURRENT_ID_AGENT int32 = 1

type AgentGroup[iAgent Agent] struct {
	world            *World
	agents           *common.SortSlice[iAgent]
	toRunPathFinding bool
	PosToAgents      map[common.Vec[int32]][]iAgent
}

func NewAgentGroup[iAgent Agent](w *World) AgentGroup[iAgent] {
	return AgentGroup[iAgent]{
		agents: common.NewSortSlice(func(a, b iAgent) int { return a.GetId() - b.GetId() }),
		world:  w,
	}
}

func (w *AgentGroup[iAgent]) GetAgentsAt(pos common.Vec[int32], condition func(iAgent) bool) (res []iAgent) {
	if condition == nil {
		condition = func(_ iAgent) bool { return true }
	}

	for i, agent := range w.PosToAgents[pos] {
		if condition(agent) {
			res = append(res, w.PosToAgents[pos][i])
		}
	}
	return res
}

func (w *AgentGroup[iAgent]) AddAgent(agent iAgent) {
	pos := agent.GetPos()
	w.agents.Insert(agent)
	agent.SetPos(pos)
	w.PosToAgents[pos] = append(w.PosToAgents[pos], agent)
}

func (w *AgentGroup[iAgent]) RemoveAgent(agent iAgent) (err error) {
	removed := w.agents.Remove(agent)
	if !removed {
		return errors.New("Error removing agent")
	}
	w.PosToAgents[agent.GetPos()] = slices.DeleteFunc(w.PosToAgents[agent.GetPos()], func(a iAgent) bool { return a.GetId() == agent.GetId() })
	return nil
}
