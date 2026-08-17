package core

import "GameOfLife/common"

type Resource string

const (
	FOOD = "FOOD"
	IRON = "IRON"
	WOOD = "WOOD"
)

type Quantity[t any, q common.Number] struct {
	What   t
	Amount q
}

var CellToResource map[CellType][]Quantity[Resource, float32] = map[CellType][]Quantity[Resource, float32]{
	WHEAT_FIELD: {{FOOD, 0.5}},
	DOCK:        {{FOOD, 1}},
	MINE:        {{IRON, 0.5}},
	FOREST:      {{WOOD, 1}},
}
