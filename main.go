package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, World!")
	size := 16
	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			g.drawRectAngle(screen, i*size, j*size, float32(size), color.White)
			g.drawLineRectAngle(screen, i*size, j*size, float32(size), color.Black, 1)

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
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello world!")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
