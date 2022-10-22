package game

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	b "go-miner/pkg/board"
	"go-miner/pkg/point"
	"golang.org/x/image/colornames"
	"image/color"
	"strconv"
	"time"
)

type GameResult int

const (
	Unknown = GameResult(iota)
	Winner
	Looser
)

type Game struct {
	Result GameResult
	IsDone bool

	width, height, mineCount, gridRadius int

	startedAt  time.Time
	finishedAt time.Time

	totalMarked int

	board      *b.Board
	boardImd   *imdraw.IMDraw
	textAtlas  *text.Atlas
	gridMatrix pixel.Matrix
	gridTexts  map[*b.Cell]*text.Text
}

func NewGame(win *pixelgl.Window, width, height, mineCount, gridRadius int, atlas *text.Atlas) *Game {
	gs := Game{
		width:      width,
		height:     height,
		mineCount:  mineCount,
		gridRadius: gridRadius,
		textAtlas:  atlas,
	}
	gs.Reset(win)
	gs.IsDone = true
	return &gs
}

func (gs *Game) Reset(win *pixelgl.Window) {
	gs.Result = Unknown
	gs.IsDone = false
	gs.totalMarked = 0
	gs.startedAt = time.Now()
	gs.finishedAt = gs.startedAt
	gs.gridTexts = map[*b.Cell]*text.Text{}
	gs.board = b.NewBoard(gs.width, gs.height, gs.mineCount)
	gs.gridMatrix = pixel.IM.Moved(win.Bounds().Center().Sub(pixel.V(
		float64(gs.width*gs.gridRadius),
		float64(gs.height*gs.gridRadius)-30,
	)))
	gs.boardImd = imdraw.New(nil)
	gs.boardImd.SetMatrix(gs.gridMatrix)
	for gp, cell := range gs.board.Field() {
		if cell.MineCount > 0 && !cell.HasBomb {
			pos := pixel.V(
				float64(gp.X*2*gs.gridRadius),
				float64(gp.Y*2*gs.gridRadius),
			)
			txt := text.New(pos, gs.textAtlas)
			txt.Color = colornames.Black
			if _, err := txt.WriteString(strconv.Itoa(cell.MineCount)); err != nil {
				panic(err)
			}
			gs.gridTexts[cell] = txt
		}
	}
}

func (gs *Game) Update(win *pixelgl.Window) {
	var selected *point.Point

	if !gs.IsDone {
		selected = gs.mousePositionToBoardPoint(win.MousePosition())
		if win.JustPressed(pixelgl.KeyEscape) {
			gs.IsDone = true
			win.Update()
		}
		if win.JustPressed(pixelgl.MouseButtonRight) {
			if cell := gs.board.At(*selected); cell != nil {
				if !cell.Opened {
					gs.board.RecursiveOpen(*selected)
					if cell.HasBomb {
						gs.finishedAt = time.Now()
						gs.Result = Looser
					}
				}
			}
		}
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			if cell := gs.board.At(*selected); cell != nil {
				if cell.Opened {
					n, ok := gs.board.TryMarkNeighbours(*selected)
					gs.totalMarked += n
					if !ok {
						gs.finishedAt = time.Now()
						gs.Result = Looser
					}
				} else {
					cell.Marked = !cell.Marked
					if cell.Marked {
						gs.totalMarked += 1
					} else {
						gs.totalMarked -= 1
					}
				}
			}
		}
	}
	if gs.WinCheck() {
		gs.board.ShowAll()
		gs.Result = Winner
	}
	if gs.Result != Unknown {
		gs.IsDone = true
	}

	gs.Draw(win, gs.boardImd, selected)
	gs.DrawUI(win)
}

func (gs *Game) WinCheck() bool {
	closedBombs := 0
	closed := 0
	for _, cell := range gs.board.Field() {
		if !cell.Opened && !cell.Marked {
			closed += 1
			if cell.HasBomb {
				closedBombs += 1
			}
		}
		if cell.HasBomb && cell.Opened {
			return false
		}
		if cell.Marked && !cell.HasBomb {
			return false
		}
	}
	return closedBombs == 0 || closedBombs == closed
}

func (gs *Game) mousePositionToBoardPoint(mp pixel.Vec) *point.Point {
	mp = gs.gridMatrix.Unproject(mp).Floor()
	gp := point.Point{
		X: int(mp.X / float64(gs.gridRadius*2)),
		Y: int(mp.Y / float64(gs.gridRadius*2)),
	}
	if mp.X < 0 {
		gp.X -= 1
	}
	if mp.Y < 0 {
		gp.Y -= 1
	}
	return &gp
}

func (gs *Game) Draw(win *pixelgl.Window, imd *imdraw.IMDraw, selected *point.Point) {
	imd.Clear()
	if gs.IsDone {
		imd.SetColorMask(colornames.Lightgray)
	} else {
		imd.SetColorMask(colornames.White)
	}
	for pos, cell := range gs.board.Field() {
		p := pixel.V(
			float64(pos.X*2*gs.gridRadius),
			float64(pos.Y*2*gs.gridRadius),
		)
		switch {
		case cell.Opened && cell.HasBomb:
			imd.Color = colornames.Red
		case cell.Opened && !cell.HasBomb:
			imd.Color = color.RGBA{R: 0xe8 + 0x10, G: 0x95 + 0x10, B: 0xd3 + 0x10, A: 0xff}
		default:
			imd.Color = color.RGBA{R: 0xe8, G: 0x95, B: 0xd3, A: 0xff}
		}
		imd.Push(
			p,
			p.Add(pixel.V(0, float64(2*gs.gridRadius))),
			p.Add(pixel.V(float64(2*gs.gridRadius), float64(2*gs.gridRadius))),
			p.Add(pixel.V(float64(2*gs.gridRadius), 0)),
		)
		imd.Polygon(0)
		if selected != nil && pos == *selected {
			imd.Color = color.RGBA{R: 0x05, G: 0x05, B: 0x05, A: 0x10}
			imd.Push(
				p,
				p.Add(pixel.V(0, float64(2*gs.gridRadius))),
				p.Add(pixel.V(float64(2*gs.gridRadius), float64(2*gs.gridRadius))),
				p.Add(pixel.V(float64(2*gs.gridRadius), 0)),
				p,
			)
			imd.Polygon(0)
		}
	}

	imd.Color = colornames.Lightgray
	imd.Push(
		pixel.V(float64(0), float64(0)),
		pixel.V(float64(gs.width*gs.gridRadius*2), float64(0)),
		pixel.V(float64(gs.width*gs.gridRadius*2), float64(gs.height*gs.gridRadius*2)),
		pixel.V(float64(0), float64(gs.height*gs.gridRadius*2)),
		pixel.V(float64(0), float64(0)),
	)
	imd.Rectangle(1)
	for i := 0; i < gs.width; i++ {
		imd.Push(
			pixel.V(float64(i*gs.gridRadius*2), float64(0)),
			pixel.V(float64(i*gs.gridRadius*2), float64(gs.height*gs.gridRadius*2)),
		)
		imd.Line(1)
	}
	for i := 0; i < gs.height; i++ {
		imd.Push(
			pixel.V(float64(0), float64(i*gs.gridRadius*2)),
			pixel.V(float64(gs.width*gs.gridRadius*2), float64(i*gs.gridRadius*2)),
		)
		imd.Line(1)
	}

	imd.Draw(win)
	imd.Clear()
	for gp, cell := range gs.board.Field() {
		if !cell.Opened && cell.Marked {
			p := pixel.V(
				float64(gp.X*2*gs.gridRadius+gs.gridRadius),
				float64(gp.Y*2*gs.gridRadius+gs.gridRadius),
			)
			imd.Color = color.RGBA{R: 0x42, G: 0xbc, B: 0xf5, A: 0xff}
			imd.Push(p)
			imd.Circle(float64(gs.gridRadius/2), 0)
		}
		if txt, ok := gs.gridTexts[cell]; ok && cell.Opened {
			txt.Draw(win, gs.gridMatrix.Moved(pixel.V(
				float64(gs.gridRadius)-txt.Bounds().W()/2,
				float64(gs.gridRadius)-txt.Bounds().H()/2,
			)))
		}
	}
	imd.Draw(win)
}

func (gs *Game) DrawUI(win *pixelgl.Window) {
	if gs.IsDone {
		return
	}
	info := text.New(pixel.ZV, gs.textAtlas)
	info.LineHeight = info.LineHeight * 1.5
	info.Color = colornames.Black

	fmt.Fprintf(info, "Отмечено %d из %d мин", gs.totalMarked, gs.mineCount)

	info.Draw(win, pixel.IM.Moved(pixel.V(
		win.Bounds().W()/2-info.Bounds().W()/2,
		win.Bounds().H()-50,
	)))

	footer := text.New(pixel.ZV, gs.textAtlas)
	footer.Color = colornames.Black
	footer.LineHeight = footer.LineHeight * 1.5

	footerStr := "Для возврата в меню нажмите Esc\n"
	footer.Dot.X -= footer.BoundsOf(footerStr).W() / 2
	footer.WriteString(footerStr)

	footerStr = "Открыть поле - правая кнопка мыши\n"
	footer.Dot.X -= footer.BoundsOf(footerStr).W() / 2
	footer.WriteString(footerStr)

	footerStr = "Поставить/Убрать флаг - левая кнопка мыши\n"
	footer.Dot.X -= footer.BoundsOf(footerStr).W() / 2
	footer.WriteString(footerStr)

	footerStr = "Клик в цифру при очевидном результате ускорит вас ;)"
	footer.Dot.X -= footer.BoundsOf(footerStr).W() / 2
	footer.WriteString(footerStr)

	footer.Draw(win, pixel.IM.Moved(pixel.V(
		win.Bounds().W()/2,
		90,
	)))
}
