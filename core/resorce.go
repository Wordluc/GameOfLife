package core

type Resource int8

const (
	FOOD Resource = iota
	IRON
)

var JobToResource map[Job][]Resource = map[Job][]Resource{
	FARMER: {FOOD},
	MINER:  {IRON},
	DOCK:   {FOOD},
}
