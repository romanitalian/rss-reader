package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.NewWithID("rss-reader")
	w := a.NewWindow("rss-reader")

	w.SetContent(widget.NewLabel("Hello World!"))
	w.ShowAndRun()
}
