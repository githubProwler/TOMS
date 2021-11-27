package colouredCircle

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type C = layout.Context
type D = layout.Dimensions

type windowState struct {
	ops          op.Ops
	startButton  widget.Clickable
	colourEditor widget.Editor
	colour       color.NRGBA
	colourString string
	theme        material.Theme
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (cc *ColouredCircle) updateColour(inputNum uint) {
	redNumber := uint8((inputNum + uint(cc.windowState.colour.R)) % 256)
	cc.windowState.colour.R = cc.windowState.colour.G
	cc.windowState.colour.G = cc.windowState.colour.B
	cc.windowState.colour.B = redNumber

	red := fmt.Sprint(cc.windowState.colour.R)
	green := fmt.Sprint(cc.windowState.colour.G)
	blue := fmt.Sprint(cc.windowState.colour.B)
	cc.windowState.colourString = "R: " + red + " G: " + green + " B: " + blue
	cc.window.Invalidate()
}

func initializeWindowState(state *windowState) {
	state.colourString = "R: 0 G: 0 B: 0"
	state.colourEditor.SingleLine = true
	state.colourEditor.Alignment = text.Middle
	state.colour = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	state.theme = *material.NewTheme(gofont.Collection())
}

func (cc *ColouredCircle) loop() error {
	var windowState windowState
	cc.windowState = &windowState
	initializeWindowState(cc.windowState)

	for {
		e := <-cc.window.Events()
		// fmt.Printf("Type: %T\n", e)

		switch e := e.(type) {
		case system.FrameEvent:
			gtx := layout.NewContext(&cc.windowState.ops, e)

			if cc.windowState.startButton.Clicked() {
				cc.fn(cc.windowState.colourEditor.Text(), cc.fnArgs)
			}

			layout.Flex{
				Axis:    layout.Vertical,
				Spacing: layout.SpaceEnd,
			}.Layout(gtx,
				layout.Rigid(
					func(gtx C) D {
						MaxX := gtx.Constraints.Max.X
						MaxY := gtx.Constraints.Max.Y
						var radius float32

						d := image.Point{X: MaxX, Y: MaxY / 2}
						radius = float32(minInt(d.X/2, d.Y/2))
						circle := clip.Circle{
							Center: f32.Point{X: float32(d.X / 2), Y: float32(d.Y / 2)},
							Radius: radius,
						}.Op(gtx.Ops)

						paint.FillShape(gtx.Ops, windowState.colour, circle)

						return layout.Dimensions{Size: d}
					},
				),
				layout.Rigid(
					func(gtx C) D {
						margin := layout.Inset{
							Top:    unit.Dp(20),
							Bottom: unit.Dp(20),
						}

						return margin.Layout(gtx,
							func(gtx C) D {
								resultColor := material.Label(&cc.windowState.theme, unit.Sp(20), windowState.colourString)
								resultColor.Alignment = text.Middle
								return resultColor.Layout(gtx)
							},
						)
					},
				),
				layout.Rigid(
					material.Editor(&cc.windowState.theme, &cc.windowState.colourEditor, "Input Color").Layout,
				),
				layout.Rigid(
					func(gtx C) D {
						margin := layout.Inset{
							Top:    unit.Dp(10),
							Bottom: unit.Dp(10),
							Right:  unit.Dp(40),
							Left:   unit.Dp(40),
						}

						return margin.Layout(gtx,
							func(gtx C) D {
								btn := material.Button(&cc.windowState.theme, &cc.windowState.startButton, "Enter")
								return btn.Layout(gtx)
							},
						)
					},
				),
			)

			e.Frame(gtx.Ops)
		case system.DestroyEvent:
			return e.Err
		}
	}
}

func (cc *ColouredCircle) createWindow(name string) {
	cc.window = app.NewWindow(
		app.Title(name),
		app.Size(unit.Dp(400), unit.Dp(400)),
	)

	if err := cc.loop(); err != nil {
		log.Fatal(err)
	}
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
