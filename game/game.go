package game

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
	CellSize    int
	selectCount int
	previousX   int
	previousY   int
	MS          minesweeperlogic.MineSweeper
}

func (g *Game) Update() error {
	x, y := ebiten.CursorPosition()
	flagPressed := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight)
	mousePressed := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	size := g.CellSize

	for i := 0; i < g.MS.FieldHeight; i++ {
		for j := 0; j < g.MS.FieldWidth; j++ {
			if j*size < x && x < (j+1)*size && i*size < y && y < (i+1)*size {
				if flagPressed {
					g.MS.Flag(j, i)
				} else if mousePressed {
					g.previousX = j
					g.previousY = i
					g.MS.Open(j, i)
					g.selectCount++
					if g.MS.Field[i][j].Count == 0 {
						g.MS.DigEmpty(j, i)
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

	size := g.CellSize

	for i := 0; i < g.MS.FieldHeight; i++ {
		for j := 0; j < g.MS.FieldWidth; j++ {

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
				if g.MS.Field[i][j].Open {
					if g.MS.Field[i][j].Bomb {
						g.drawRectAngle(screen, j*size, i*size, float32(size), color.RGBA{0xff, 0x00, 0x00, 0xff})
						g.drawLineRectAngle(screen, j*size, i*size, float32(size), color.White, 1)

					} else {
						g.drawRectAngle(screen, j*size, i*size, float32(size), color.Black)
						g.drawLineRectAngle(screen, j*size, i*size, float32(size), color.White, 1)
					}
					g.drawText(screen, strconv.Itoa(g.MS.Field[i][j].Count), j*size, i*size, float64(g.CellSize/2), color.White)
					// ebitenutil.DebugPrintAt(screen, strconv.Itoa(g.MS.Field[i][j].Count), j*size, i*size)
				} else if g.MS.Field[i][j].Flag {
					g.drawRectAngle(screen, j*size, i*size, float32(size), color.RGBA{0xff, 0xff, 0x00, 0xff})
					g.drawLineRectAngle(screen, j*size, i*size, float32(size), color.White, 1)
				} else {
					g.drawLineRectAngle(screen, j*size, i*size, float32(size), color.White, 1)
				}
			}

		}
	}

	g.MS.IsGameClear()

	if g.MS.GameOver {
		g.drawText(screen, "Game Over!", 0, 0, 30, color.White)
		g.MS.AllOpen()
		if g.selectCount == 1 {
			g.MS.Init(g.MS.FieldWidth, g.MS.FieldHeight)
			g.MS.SummonBomb()
			g.MS.CountBomb()
			g.selectCount = 0
		}
	}

	if g.MS.GameClear {
		g.drawText(screen, "Game Clear!", 0, 0, 30, color.White)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.CellSize * g.MS.FieldWidth, g.CellSize * g.MS.FieldHeight
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
