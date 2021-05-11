package main

import (
	"log"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type field1 struct {
	app.Compo
	Label        string
	Value        string
	Help         string
	HelpClass    string
	ReportChange func(v string) (string, error)
}

func newField1(v string, l string) *field1 {
	return &field1{
		Label: l,
		Value: v,
		Help:  "All good",
	}
}

func (f *field1) WithReportChange(fn func(v string) (string, error)) *field1 {
	f.ReportChange = fn
	return f
}

func (f *field1) Render() app.UI {
	return app.Div().Body(
		app.Label().Text(f.Label).Style("color", "blue"),
		app.Input().Type("text").Value(f.Value).OnInput(f.valueChanged, f.Label),
		app.P().Style("color", f.HelpClass).Body(
			app.Text("Help: "),
			app.Text(f.Help),
		),
	)
}

func (f *field1) valueChanged(ctx app.Context, e app.Event) {
	f.Value = ctx.JSSrc().Get("value").String()
	if f.ReportChange != nil {
		v, err := f.ReportChange(f.Value)
		if err == nil {
			f.HelpClass = "green"
		} else {
			f.HelpClass = "red"
		}
		f.Help = v
		log.Printf("ReportChange Help: %s, Class:%s", f.Help, f.HelpClass)
	}
}
