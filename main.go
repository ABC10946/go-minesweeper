package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/ABC10946/minesweeper/minesweeperlogic"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	ms minesweeperlogic.MineSweeper
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	x, y := ebiten.CursorPosition()
	mousePressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	size := 16
	for i := 0; i < g.ms.FieldHeight; i++ {
		for j := 0; j < g.ms.FieldWidth; j++ {

			if i*size < x && x < (i+1)*size && j*size < y && y < (j+1)*size {
				if mousePressed {
					g.drawRectAngle(screen, i*size, j*size, float32(size), color.RGBA{0x00, 0xff, 0x00, 0xff})
					g.drawLineRectAngle(screen, i*size, j*size, float32(size), color.Black, 1)

				} else {
					g.drawRectAngle(screen, i*size, j*size, float32(size), color.RGBA{0xff, 0xff, 0x00, 0xff})
					g.drawLineRectAngle(screen, i*size, j*size, float32(size), color.Black, 1)
				}
			} else {
				if g.ms.Field[j][i].Bomb {
					g.drawRectAngle(screen, i*size, j*size, float32(size), color.RGBA{0xff, 0x00, 0x00, 0xff})

				} else {
					g.drawRectAngle(screen, i*size, j*size, float32(size), color.White)
					g.drawLineRectAngle(screen, i*size, j*size, float32(size), color.Black, 1)
				}
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
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
	ebiten.SetWindowSize(640, 640)
	ebiten.SetWindowTitle("Hello world!")
	game := Game{}
	game.ms.Init(20, 20)
	game.ms.SummonBomb()
	fmt.Println(game.ms.Field)
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
