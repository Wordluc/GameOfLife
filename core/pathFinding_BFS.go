package core

import (
	"GameOfLife/common"
	"math"
	"slices"
)

func PerformPathFinding_BFS(m *Map[BaseCell], starts []common.Vec[int32]) (res Map[int16], err error) {
	var queue *common.Queue[Quantity[common.Vec[int32], int16]] = common.NewQueue[Quantity[common.Vec[int32], int16]](nil, nil)
	var visited map[common.Vec[int32]]bool = make(map[common.Vec[int32]]bool, m.size.X*m.size.Y)
	res = NewMap[int16](m.size)
	res.SetEachElement(func(x, y int32) *int16 {
		if n, _ := m.GetCell(common.Vec[int32]{X: x, Y: y}); slices.Contains([]CellType{STONE, WATER}, n.cellType) {
			return new(int16(math.MaxInt8) + 1)
		}
		return new(int16(math.MaxInt8))
	})
	for i := range starts {
		v, err := res.GetCell(starts[i])
		if err != nil {
			return res, err
		}
		v = new(int16(0))
		res.SetRawCell(v, starts[i])
		q := Quantity[common.Vec[int32], int16]{
			What:   starts[i],
			Amount: -1,
		}
		queue.Enqueue(q)
	}
	for {
		toSee, end := queue.Denqueue()
		if end {
			break
		}
		if visited[toSee.What] {
			continue
		}
		cost := int16(toSee.Amount + 1)
		if c, _ := m.GetCell(toSee.What); slices.Contains([]CellType{STONE, WATER}, c.cellType) {
			continue
		}
		res.SetRawCell(&cost, toSee.What)
		visited[toSee.What] = true
		n, _ := res.GetNeighborhoodCells(toSee.What, common.Vec[int32]{X: 3, Y: 3})
		for i := range n {
			if i.IsEqual(toSee.What) {
				continue
			}
			if visited[i] {
				continue
			}
			queue.Enqueue(Quantity[common.Vec[int32], int16]{What: i, Amount: cost})
		}
	}

	return res, nil
}
