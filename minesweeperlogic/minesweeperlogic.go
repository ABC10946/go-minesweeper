package minesweeperlogic

import (
	"log"
	"math/rand"
)

type Cell struct {
	Open  bool
	Count int
	Bomb  bool
	Flag  bool
}

type MineSweeper struct {
	FieldWidth  int
	FieldHeight int
	Field       [][]Cell
	GameOver    bool
	TotalBomb   int
	GameClear   bool
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
	totalBomb := 0
	if ms.Field == nil {
		log.Fatal("Error: field is empty")
	}

	for i := 0; i < ms.FieldHeight; i++ {
		for j := 0; j < ms.FieldWidth; j++ {
			if rand.NormFloat64() > 0.9 {
				ms.Field[i][j].Bomb = true
				totalBomb++
			}
		}
	}

	ms.TotalBomb = totalBomb
}

func (ms *MineSweeper) cellCountBomb(x, y int) int {
	if ms.Field[y][x].Bomb {
		return -1
	}

	// x, y
	neighborIdx := [][]int{
		{-1, -1},
		{0, -1},
		{1, -1},
		{-1, 0},
		{1, 0},
		{-1, 1},
		{0, 1},
		{1, 1},
	}

	count := 0

	for i := 0; i < 8; i++ {
		xshift := neighborIdx[i][0]
		yshift := neighborIdx[i][1]
		neighborCellX := x + xshift
		neighborCellY := y + yshift

		if !ms.isOverWall(neighborCellX, neighborCellY) {
			if ms.Field[neighborCellY][neighborCellX].Bomb {
				count++
			}
		}
	}

	return count
}

func (ms *MineSweeper) CountBomb() {
	for i := 0; i < ms.FieldHeight; i++ {
		for j := 0; j < ms.FieldWidth; j++ {
			ms.Field[i][j].Count = ms.cellCountBomb(j, i)
		}
	}
}

func (ms *MineSweeper) isOverWall(x, y int) bool {
	if x < 0 || y < 0 {
		return true
	}

	if x >= ms.FieldWidth {
		return true
	}

	if y >= ms.FieldHeight {
		return true
	}

	return false
}

func (ms *MineSweeper) Open(x, y int) {
	if !ms.Field[y][x].Flag {
		if ms.Field[y][x].Bomb {
			ms.GameOver = true
		}

		if !ms.Field[y][x].Bomb {
			ms.Field[y][x].Open = true
		}
	}
}

func (ms *MineSweeper) AllOpen() {
	for i := 0; i < ms.FieldHeight; i++ {
		for j := 0; j < ms.FieldWidth; j++ {
			ms.Field[i][j].Open = true
		}
	}
}

func (ms *MineSweeper) Flag(x, y int) {
	ms.Field[y][x].Flag = !ms.Field[y][x].Flag
}

func (ms *MineSweeper) IsGameClear() {
	emptyCellCount := ms.FieldHeight*ms.FieldWidth - ms.TotalBomb
	openedCellCount := 0

	for i := 0; i < ms.FieldHeight; i++ {
		for j := 0; j < ms.FieldWidth; j++ {
			if ms.Field[i][j].Open {
				openedCellCount++
			}
		}
	}
	if openedCellCount == emptyCellCount {
		ms.GameClear = true
	}
}

// カウント0のセルを見つけたとき、周囲に0のセルがあったらそれも掘り出す
func (ms *MineSweeper) DigEmpty(x, y int) {
	neighborIdx := [][]int{
		{-1, -1},
		{0, -1},
		{1, -1},
		{-1, 0},
		{1, 0},
		{-1, 1},
		{0, 1},
		{1, 1},
	}

	for i := 0; i < 8; i++ {
		xshift := neighborIdx[i][0]
		yshift := neighborIdx[i][1]
		neighborCellX := x + xshift
		neighborCellY := y + yshift

		if !ms.isOverWall(neighborCellX, neighborCellY) {
			if !ms.Field[neighborCellY][neighborCellX].Open {
				ms.Field[neighborCellY][neighborCellX].Open = true
				if ms.Field[neighborCellY][neighborCellX].Count == 0 {
					ms.DigEmpty(neighborCellX, neighborCellY)
				}
			}
		}
	}
}
