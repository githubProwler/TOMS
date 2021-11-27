package colouredCircle

import (
	"image"
	"image/color"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/event"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type C = layout.Context
type D = layout.Dimensions

type windowState struct {
	ops         op.Ops
	startButton widget.Clickable
	started     bool
	colour      color.NRGBA
	theme       material.Theme
	c           chan bool
	closed      bool
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func timer(c chan bool) {
	for {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(2000)+100))
		c <- true
	}
}

func (cc *ColouredCircle) updateColour(inputNum uint) {
	redNumber := uint8((inputNum + uint(cc.windowState.colour.R)) % 256)
	cc.windowState.colour.R = cc.windowState.colour.G
	cc.windowState.colour.G = cc.windowState.colour.B
	cc.windowState.colour.B = redNumber

	cc.window.Invalidate()
}

func initializeWindowState(state *windowState) {
	state.colour = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	state.theme = *material.NewTheme(gofont.Collection())
	state.c = make(chan bool)
	go timer(state.c)
}

func (cc *ColouredCircle) handleWindowEvent(e event.Event) error {
	// log.Printf("Type: %T\n", e)

	switch e := e.(type) {
	case system.FrameEvent:
		gtx := layout.NewContext(&cc.windowState.ops, e)

		if cc.windowState.startButton.Clicked() {
			cc.windowState.started = !cc.windowState.started
		}

		layout.Flex{
			Axis:    layout.Vertical,
			Spacing: layout.SpaceEnd,
		}.Layout(gtx,
			layout.Rigid(
				func(gtx C) D {
					MaxX := gtx.Constraints.Max.X
					MaxY := gtx.Constraints.Max.Y

					d := image.Point{X: MaxX, Y: MaxY - 37}
					// radius = float32(minInt(d.X/2, d.Y/2))
					rect := image.Rect(0, 0, d.X, d.Y)
					r := clip.Rect(rect).Op()
					paint.FillShape(gtx.Ops, cc.windowState.colour, r)

					return layout.Dimensions{Size: d}
				},
			),
			layout.Rigid(
				func(gtx C) D {
					var text string
					if cc.windowState.started {
						text = "Stop"
					} else {
						text = "Start"
					}
					btn := material.Button(&cc.windowState.theme, &cc.windowState.startButton, text)
					return btn.Layout(gtx)
				},
			),
		)

		e.Frame(gtx.Ops)
	case system.DestroyEvent:
		cc.windowState.closed = true
		return e.Err
	}

	return nil
}

func (cc *ColouredCircle) loop() error {
	var windowState windowState
	cc.windowState = &windowState
	initializeWindowState(cc.windowState)

	for {
		select {
		case e := <-cc.window.Events():
			err := cc.handleWindowEvent(e)
			if err != nil || cc.windowState.closed {
				return err
			}
		case <-cc.windowState.c:
			if cc.windowState.started {
				cc.fn(strconv.Itoa(rand.Intn(256)), cc.fnArgs)
			}
		}

	}
}

func (cc *ColouredCircle) createWindow(name string) {
	cc.window = app.NewWindow(
		app.Title(name),
		app.Size(unit.Dp(400), unit.Dp(400)),
	)

	if err := cc.loop(); err != nil {
		log.Print("Hello")
		log.Fatal(err)
	}
	log.Print("Hello")
	os.Exit(0)
}

type OnClickFn func(string, interface{})

type ColouredCircle struct {
	fn          OnClickFn
	fnArgs      interface{}
	window      *app.Window
	windowState *windowState
}

func (cc *ColouredCircle) AddColor(inputNum int) {
	inputNum = inputNum % 256
	if inputNum < 0 {
		inputNum += 256
	}
	cc.updateColour(uint(inputNum))
}

// Main must be called last from the program main function.
// On most platforms Main blocks forever, for Android and
// iOS it returns immediately to give control of the main
// thread back to the system.
//
// Calling Main is necessary because some operating systems
// require control of the main thread of the program for
// running windows.
func (cc *ColouredCircle) Main(name string, fn OnClickFn, fnArgs interface{}) {
	cc.fn = fn
	cc.fnArgs = fnArgs
	go cc.createWindow(name)
	app.Main()
}
