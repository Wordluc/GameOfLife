package core

import "GameOfLife/common"

var JobToCells = map[Job]CellType{
	FARMER:     WHEAT_FIELD,
	MINER:      MINE,
	WOODCUTTER: FOREST,
	FISHERMAN:  DOCK,
}

var CellToJobs = common.ReverseMapToValue(JobToCells)

var JobToConsumingCost = map[Job][]Quantity[Resource, float32]{
	FARMER:     {{FOOD, 1}},
	MINER:      {{FOOD, 1}},
	WOODCUTTER: {{FOOD, 1}},
	FISHERMAN:  {{FOOD, 1.5}},
}

var CellTypeToResource map[CellType][]Quantity[Resource, float32] = map[CellType][]Quantity[Resource, float32]{
	WHEAT_FIELD: {{FOOD, 2}},
	DOCK:        {{FOOD, 3}},
	MINE:        {{IRON, 0.5}},
	FOREST:      {{WOOD, 1}},
}
