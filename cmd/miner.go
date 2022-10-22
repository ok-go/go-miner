package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"go-miner/pkg/game"
	"go-miner/pkg/utils"
	"golang.org/x/image/colornames"
	"math/rand"
	"time"
	"unicode"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	f, err := utils.LoadTTF("AverageMono.ttf", 14)
	if err != nil {
		panic(err)
	}

	baseAtlas = text.NewAtlas(
		f,
		text.ASCII,
		text.RangeTable(unicode.Cyrillic),
	)

	pixelgl.Run(func() {
		win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
			Title:  "Miner",
			Bounds: pixel.R(0, 0, 800, 600),
		})
		if err != nil {
			panic(err)
		}
		defer win.Destroy()

		loop(win)
	})
}

var baseAtlas *text.Atlas
var gs *game.Game

func loop(win *pixelgl.Window) {
	gs = game.NewGame(win, 10, 10, 20, 20, baseAtlas)

	for !win.Closed() {
		win.Clear(colornames.Whitesmoke)
		gs.Update(win)
		if gs.IsDone {
			MenuLoop(win)
		}
		win.Update()
	}
}

func MenuLoop(win *pixelgl.Window) {
	if win.JustPressed(pixelgl.KeyEscape) {
		win.SetClosed(true)
		return
	}

	resultText := text.New(pixel.ZV, baseAtlas)
	resultText.Color = colornames.Black
	resultText.LineHeight = resultText.LineHeight * 1.4
	switch gs.Result {
	case game.Looser:
		fmt.Fprintf(resultText, "Случайно подорвался? В следующий раз повезет больше!")
	case game.Winner:
		fmt.Fprintf(resultText, "Поздравляю! Справишься ли ты с новым полем?")
	default:
		fmt.Fprintf(resultText, "Розовый - хит сезона")
	}
	resultText.Draw(win, pixel.IM.Moved(pixel.V(
		win.Bounds().W()/2-resultText.Bounds().W()/2,
		win.Bounds().H()-50)))

	pressAnyKeyTitle := text.New(pixel.ZV, baseAtlas)
	pressAnyKeyTitle.Color = colornames.Black
	pressAnyKeyTitle.LineHeight = pressAnyKeyTitle.LineHeight * 1.4
	WriteText(pressAnyKeyTitle, "Нажмите Esc для выхода или любую другую для начала игры\n")
	pressAnyKeyTitle.Draw(win, pixel.IM.Moved(pixel.V(
		win.Bounds().W()/2-pressAnyKeyTitle.Bounds().W()/2,
		50)))

	if s := win.Typed(); s != "" {
		gs.Reset(win)
	}
}

func WriteText(to *text.Text, msg string) {
	if _, err := to.WriteString(msg); err != nil {
		panic(err)
	}
}
