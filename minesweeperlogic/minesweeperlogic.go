package minesweeperlogic

import "math/rand"

type Cell struct {
	Open  bool
	Count int
	Bomb  bool
}

type MineSweeper struct {
	FieldWidth  int
	FieldHeight int
	Field       [][]Cell
}

func (ms *MineSweeper) Init(fieldWidth int, fieldHeight int) {
	var field [][]Cell
	ms.FieldWidth = fieldWidth
	ms.FieldHeight = fieldHeight

	for i := 0; i < fieldHeight; i++ {
		field = append(field, make([]Cell, fieldWidth))
	}
	ms.Field = field
}

func (ms *MineSweeper) SummonBomb() {
	if ms.Field == nil {
		panic("Error: field is empty")
	}

	fieldHeight := len(ms.Field)
	fieldWidth := len(ms.Field[0])

	for i := 0; i < fieldHeight; i++ {
		for j := 0; j < fieldWidth; j++ {
			if rand.NormFloat64() > 0.9 {
				ms.Field[j][i].Bomb = true
			}
		}
	}
}
