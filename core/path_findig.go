package core

import (
	"GameOfLife/common"
	"slices"
	"sort"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type pathFindigInformation struct {
	pos             common.Vec[int32]
	origin          *pathFindigInformation
	weightFromStart int32
	weightToGoal    int32
}

func (p *pathFindigInformation) retriavePath() (res []common.Vec[int32]) {
	var t *pathFindigInformation
	t = p
	for {
		if t == nil {
			slices.Reverse(res)
			return res
		}
		res = append(res, t.pos)
		t = t.origin
	}

}
func Search(s []*pathFindigInformation, a common.Vec[int32]) int {
	for i := range s {
		if s[i].pos.IsEqual(a) {
			return i
		}
	}
	return -1
}
func InsertSorted(s []*pathFindigInformation, e *pathFindigInformation) []*pathFindigInformation {
	ef := e.weightFromStart + e.weightToGoal

	i := sort.Search(len(s), func(i int) bool {
		f := s[i].weightFromStart + s[i].weightToGoal

		if f != ef {
			return f > ef
		}

		return s[i].weightToGoal > e.weightToGoal
	})

	s = append(s, nil)
	copy(s[i+1:], s[i:])
	s[i] = e

	return s
}

func fromAtoB(a, b common.Vec[int32]) int32 {
	return common.DistanceAtoBVecShev(a, b)
}

func foreachNeighboarhood(m *Map, pos common.Vec[int32], callback func(x, y int32) (stop bool)) {
	for iy := range int32(3) {
		for ix := range int32(3) {
			worldX := pos.X + (ix - 1)
			worldY := pos.Y + (iy - 1)

			if worldX == pos.X && worldY == pos.Y {
				continue
			}

			if worldX < 0 || worldX >= m.size.X {
				continue
			}

			if worldY < 0 || worldY >= m.size.Y {
				continue
			}

			if callback(worldX, worldY) {
				return
			}
		}
	}
}
func drawDebugCell(pos common.Vec[int32], c rl.Color) {
	rl.BeginDrawing()
	rl.DrawRectangle(pos.X*10, pos.Y*10, 10, 10, c)
	rl.EndDrawing()
}

func PerformPathFindig(m *Map, start, goal common.Vec[int32]) []common.Vec[int32] {
	var discovered map[common.Vec[int32]]*pathFindigInformation = make(map[common.Vec[int32]]*pathFindigInformation)
	var toDiscover []*pathFindigInformation = make([]*pathFindigInformation, 0)
	startingOrigin := pathFindigInformation{
		pos: start,
	}
	discoverNeighboardhood := func(origin pathFindigInformation) func(x, y int32) (stop bool) {
		return func(x, y int32) (stop bool) {
			pos := common.Vec[int32]{X: x, Y: y}
			if c, err := m.GetCell(pos); err == nil {
				if slices.Contains([]CellType{WATER}, c.blockType) {
					return false
				}
			}

			weightFromStart := origin.weightFromStart + 2
			if !slices.Contains([]float32{0, 90, 180, 270, 360}, pos.Clone().Sub(goal).Angle()) {
				weightFromStart += 1
			}
			if pos.IsEqual(goal) {
				toDiscover = InsertSorted(toDiscover, &pathFindigInformation{
					pos:             pos,
					origin:          &origin,
					weightFromStart: weightFromStart,
					weightToGoal:    fromAtoB(pos, goal),
				})
				return true
			}
			weightToGoal := fromAtoB(pos, goal)
			if _, ok := discovered[pos]; ok {
				return
			}
			if id := Search(toDiscover, pos); id != -1 {
				a := toDiscover[id]
				if weightToGoal+weightFromStart < a.weightToGoal+a.weightFromStart {
					a.origin = &origin
					a.weightFromStart = weightFromStart
					a.weightToGoal = weightToGoal
					toDiscover = slices.Delete(toDiscover, id, id+1)
					toDiscover = InsertSorted(toDiscover, a)
				}
				return
			}
			toDiscover = InsertSorted(toDiscover, &pathFindigInformation{
				pos:             pos,
				origin:          &origin,
				weightFromStart: weightFromStart,
				weightToGoal:    fromAtoB(pos, goal),
			})
			return false
		}
	}
	foreachNeighboarhood(m, start, discoverNeighboardhood(startingOrigin))
	for {
		if len(toDiscover) == 0 {
			return nil
		}
		explored := toDiscover[0]
		toDiscover = slices.Delete(toDiscover, 0, 1)

		if explored.pos.IsEqual(goal) {
			return explored.retriavePath()
		}
		discovered[explored.pos] = explored
		foreachNeighboarhood(m, explored.pos, discoverNeighboardhood(*explored))

	}

}
