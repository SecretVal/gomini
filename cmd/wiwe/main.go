package main

import (
	"crypto/tls"
	"fmt"
	"image/color"
	"log"
	"os"
	"strings"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
)

var TLSConfig = tls.Config{
	InsecureSkipVerify: true,
}

var PORT int = 1965
var PREFIX = "gemini://"

func read_response(conn *tls.Conn) string {
	var sb strings.Builder
	buf := make([]byte, 1024)

	i, _ := conn.Read(buf)
	for i > 0 {
		sb.Write(buf)
		buf = make([]byte, 1024)
		i, _ = conn.Read(buf)
	}
	return sb.String()
}

func host_from_string(url string) string {
	return strings.Split(strings.TrimPrefix(url, PREFIX), "/")[0]
}

func make_gemini_query(url string) string {
	host := host_from_string(url)
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", host, PORT), &TLSConfig)
	if err != nil {
		panic(err)
	}

	url_buf := []byte(fmt.Sprintf("%s\r\n", url))

	_, err = conn.Write(url_buf)
	if err != nil {
		panic(err)
	}
	buf := read_response(conn)

	return buf
}

func main() {
	prog := os.Args[0]
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Printf("Usage: <%s> <url>\n", prog)
		fmt.Printf("ERROR: Did not specify url\n")
		os.Exit(1)
	}
	url := args[0]

	buf := make_gemini_query(url)
	log.Printf("Received %d bytes", len(buf))

	go func() {
		window := new(app.Window)
		err := render_window(window, buf)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func render_window(window *app.Window, str string) error {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	th.Bg = color.NRGBA{R: 18, G: 18, B: 18, A: 255}
	var ops op.Ops
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			// This graphics context is used for managing the rendering state.
			gtx := app.NewContext(&ops, e)

			// Define an large label with an appropriate text:
			body := material.Label(th, 16, str)

			// Change the color of the label.
			black := color.NRGBA{R: 0, G: 0, B: 0, A: 255}
			body.Color = black

			// Draw the label to the graphics context.
			body.Layout(gtx)

			// Pass the drawing operations to the GPU.
			e.Frame(gtx.Ops)
		}
	}
}
