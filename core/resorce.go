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
