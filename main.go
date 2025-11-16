package main

import (
	"fmt"
	"image/color"
	"log"
	"strconv"

	"github.com/ABC10946/minesweeper/minesweeperlogic"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	cellSize int
	ms       minesweeperlogic.MineSweeper
}

func (g *Game) Update() error {
	x, y := ebiten.CursorPosition()
	flagPressed := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight)
	size := g.cellSize

	for i := 0; i < g.ms.FieldHeight; i++ {
		for j := 0; j < g.ms.FieldWidth; j++ {
			if i*size < x && x < (i+1)*size && j*size < y && y < (j+1)*size {
				if flagPressed {
					g.ms.Flag(i, j)
					fmt.Printf("FLAG PRESSED %v\n", g.ms.Field[j][i].Flag)
				}
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	x, y := ebiten.CursorPosition()
	mousePressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	flagPressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)

	size := g.cellSize

	for i := 0; i < g.ms.FieldHeight; i++ {
		for j := 0; j < g.ms.FieldWidth; j++ {

			if i*size < x && x < (i+1)*size && j*size < y && y < (j+1)*size {
				if mousePressed {
					g.drawRectAngle(screen, i*size, j*size, float32(size), color.RGBA{0x00, 0xff, 0x00, 0xff})
					g.drawLineRectAngle(screen, i*size, j*size, float32(size), color.White, 1)
					g.ms.Open(i, j)
					if g.ms.Field[j][i].Count == 0 {
						g.ms.DigEmpty(i, j)
					}
				} else if flagPressed {
					g.drawRectAngle(screen, i*size, j*size, float32(size), color.RGBA{0xff, 0xff, 0x00, 0xff})
					g.drawLineRectAngle(screen, i*size, j*size, float32(size), color.White, 1)
				} else {
					g.drawRectAngle(screen, i*size, j*size, float32(size), color.RGBA{0xff, 0x00, 0xff, 0xff})
					g.drawLineRectAngle(screen, i*size, j*size, float32(size), color.White, 1)
				}
			} else {
				// normal display process
				if g.ms.Field[j][i].Open {
					if g.ms.Field[j][i].Bomb {
						g.drawRectAngle(screen, i*size, j*size, float32(size), color.RGBA{0xff, 0x00, 0x00, 0xff})
						g.drawLineRectAngle(screen, i*size, j*size, float32(size), color.White, 1)

					} else {
						g.drawRectAngle(screen, i*size, j*size, float32(size), color.Black)
						g.drawLineRectAngle(screen, i*size, j*size, float32(size), color.White, 1)
					}
					ebitenutil.DebugPrintAt(screen, strconv.Itoa(g.ms.Field[j][i].Count), i*size, j*size)
				} else if g.ms.Field[j][i].Flag {
					g.drawRectAngle(screen, i*size, j*size, float32(size), color.RGBA{0xff, 0xff, 0x00, 0xff})
					g.drawLineRectAngle(screen, i*size, j*size, float32(size), color.White, 1)
				} else {
					g.drawLineRectAngle(screen, i*size, j*size, float32(size), color.White, 1)
				}
			}

		}
	}

	g.ms.IsGameClear()

	if g.ms.GameOver {
		ebitenutil.DebugPrint(screen, "Game Over!")
		g.ms.AllOpen()
	}

	if g.ms.GameClear {
		ebitenutil.DebugPrint(screen, "Game Clear!")
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 320
}

func (g *Game) drawRectAngle(screen *ebiten.Image, x, y int, size float32, colorData color.Color) {
	var path vector.Path
	path.MoveTo(0, 0)
	path.LineTo(0, 1*size)
	path.LineTo(1*size, 1*size)
	path.LineTo(1*size, 0)
	path.Close()

	var newPath vector.Path
	op := &vector.AddPathOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	newPath.AddPath(&path, op)

	drawOp := &vector.DrawPathOptions{}
	drawOp.AntiAlias = false
	drawOp.ColorScale.ScaleWithColor(colorData)

	vector.FillPath(screen, &newPath, nil, drawOp)
}

func (g *Game) drawLineRectAngle(screen *ebiten.Image, x, y int, size float32, colorData color.Color, lineWidth float32) {
	var path vector.Path
	path.MoveTo(0, 0)
	path.LineTo(0, 1*size)
	path.LineTo(1*size, 1*size)
	path.LineTo(1*size, 0)
	path.Close()

	var newPath vector.Path
	op := &vector.AddPathOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	newPath.AddPath(&path, op)

	drawOp := &vector.DrawPathOptions{}
	drawOp.AntiAlias = false
	drawOp.ColorScale.ScaleWithColor(colorData)

	strokeOp := &vector.StrokeOptions{}
	strokeOp.Width = lineWidth
	strokeOp.LineJoin = vector.LineJoinRound
	vector.StrokePath(screen, &newPath, strokeOp, drawOp)
}

func main() {
	ebiten.SetWindowSize(1280, 1280)
	ebiten.SetWindowTitle("MINESWEEPER")
	game := Game{}
	game.cellSize = 16
	game.ms.Init(20, 20)
	game.ms.SummonBomb()
	game.ms.CountBomb()
	fmt.Println(game.ms.Field)
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
