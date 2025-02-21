package main

import (
	"fmt"
	"os"

	gemini "github.com/secretval/wiwe/cmd/wiwe/protocols/gemini"

    "log"
    "image"
    "strings"
    "image/color"

    "gioui.org/app"
    "gioui.org/op"
    "gioui.org/layout"
    "gioui.org/text"
    "gioui.org/font/gofont"
    "gioui.org/op/paint"
    "gioui.org/op/clip"
    "gioui.org/widget/material"
)

func main() {
	prog := os.Args[0]
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Printf("Usage: <%s> <url>\n", prog)
		fmt.Printf("ERROR: Did not specify url\n")
		os.Exit(1)
	}
	url := args[0]

	req, err := gemini.ParseGeminiRequest(url, gemini.PORT)
	if err != true {
		fmt.Print("ERROR: Invalid URL")
		os.Exit(1)
	}

	res := gemini.MakeGeminiQuery(req)
    log.Printf("Received %d bytes", len(res.Body))

    go func() {
        w := new(app.Window)
        err := display(w, res.Body)
        if err != nil {
            log.Fatal(err)
        }
        os.Exit(0)
    }()
    app.Main()
}


var list layout.List

func display(w *app.Window, buf string) error {
    theme := material.NewTheme()
    theme.Bg = color.NRGBA{18,18,18,255}
    theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
    var ops op.Ops
    for {
        ev := w.Event()
        switch e := ev.(type) {
        case app.DestroyEvent:
            return e.Err
            case app.FrameEvent: 
            gtx := app.NewContext(&ops, e)

            DrawRect(clip.Rect{Max: image.Pt(gtx.Constraints.Max.X,gtx.Constraints.Max.Y)}, theme.Bg, &ops)

            list.Axis = layout.Vertical
            lines := strings.Split(buf, "\n")
            list.Layout(gtx, len(lines), func(gtx layout.Context, i int) layout.Dimensions {
                line := lines[i]
                label := material.Label(theme, 16, fmt.Sprintf("%s", line))
                label.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
                if strings.HasPrefix(line, "=>") {
                    label.Color = color.NRGBA{R: 125, G: 125, B: 255, A: 255}
                }
                return label.Layout(gtx)
            })

            e.Frame(gtx.Ops)
        }
    }

}
func DrawRect(rect clip.Rect, color color.NRGBA, ops *op.Ops)  {
    defer rect.Push(ops).Pop()
    paint.ColorOp{Color: color}.Add(ops)
    paint.PaintOp{}.Add(ops)
}
