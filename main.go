package main

import (
	"bytes"
	"image/color"
	"log"
	"strconv"

	"github.com/ABC10946/minesweeper/minesweeperlogic"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	mplusFaceSource *text.GoTextFaceSource
)

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}

	mplusFaceSource = s

}

type Game struct {
	cellSize    int
	selectCount int
	previousX   int
	previousY   int
	ms          minesweeperlogic.MineSweeper
}

func (g *Game) Update() error {
	x, y := ebiten.CursorPosition()
	flagPressed := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight)
	mousePressed := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	size := g.cellSize

	for i := 0; i < g.ms.FieldHeight; i++ {
		for j := 0; j < g.ms.FieldWidth; j++ {
			if j*size < x && x < (j+1)*size && i*size < y && y < (i+1)*size {
				if flagPressed {
					g.ms.Flag(j, i)
				} else if mousePressed {
					g.previousX = j
					g.previousY = i
					g.ms.Open(j, i)
					g.selectCount++
					if g.ms.Field[i][j].Count == 0 {
						g.ms.DigEmpty(j, i)
					}
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

			if j*size < x && x < (j+1)*size && i*size < y && y < (i+1)*size {
				if mousePressed {
					g.drawRectAngle(screen, j*size, i*size, float32(size), color.RGBA{0x00, 0xff, 0x00, 0xff})
					g.drawLineRectAngle(screen, j*size, i*size, float32(size), color.White, 1)
				} else if flagPressed {
					g.drawRectAngle(screen, j*size, i*size, float32(size), color.RGBA{0xff, 0xff, 0x00, 0xff})
					g.drawLineRectAngle(screen, j*size, i*size, float32(size), color.White, 1)
				} else {
					g.drawRectAngle(screen, j*size, i*size, float32(size), color.RGBA{0xff, 0x00, 0xff, 0xff})
					g.drawLineRectAngle(screen, j*size, i*size, float32(size), color.White, 1)
				}
			} else {
				// normal display process
				if g.ms.Field[i][j].Open {
					if g.ms.Field[i][j].Bomb {
						g.drawRectAngle(screen, j*size, i*size, float32(size), color.RGBA{0xff, 0x00, 0x00, 0xff})
						g.drawLineRectAngle(screen, j*size, i*size, float32(size), color.White, 1)

					} else {
						g.drawRectAngle(screen, j*size, i*size, float32(size), color.Black)
						g.drawLineRectAngle(screen, j*size, i*size, float32(size), color.White, 1)
					}
					g.drawText(screen, strconv.Itoa(g.ms.Field[i][j].Count), j*size, i*size, float64(g.cellSize/2), color.White)
					// ebitenutil.DebugPrintAt(screen, strconv.Itoa(g.ms.Field[i][j].Count), j*size, i*size)
				} else if g.ms.Field[i][j].Flag {
					g.drawRectAngle(screen, j*size, i*size, float32(size), color.RGBA{0xff, 0xff, 0x00, 0xff})
					g.drawLineRectAngle(screen, j*size, i*size, float32(size), color.White, 1)
				} else {
					g.drawLineRectAngle(screen, j*size, i*size, float32(size), color.White, 1)
				}
			}

		}
	}

	g.ms.IsGameClear()

	if g.ms.GameOver {
		g.drawText(screen, "Game Over!", 0, 0, 30, color.White)
		g.ms.AllOpen()
		if g.selectCount == 1 {
			g.ms.Init(g.ms.FieldWidth, g.ms.FieldHeight)
			g.ms.SummonBomb()
			g.ms.CountBomb()
			g.selectCount = 0
		}
	}

	if g.ms.GameClear {
		g.drawText(screen, "Game Clear!", 0, 0, 30, color.White)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.cellSize * g.ms.FieldWidth, g.cellSize * g.ms.FieldHeight
}

func (g *Game) drawText(screen *ebiten.Image, printText string, x, y int, size float64, colorData color.Color) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleWithColor(colorData)
	text.Draw(screen, printText, &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   size,
	}, op)

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
	game := Game{}
	game.cellSize = 70
	width := 10
	height := 20
	game.ms.Init(width, height)
	game.ms.SummonBomb()
	game.ms.CountBomb()
	ebiten.SetWindowSize(game.cellSize*width, game.cellSize*height)
	ebiten.SetWindowTitle("MINESWEEPER")
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
