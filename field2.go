package main

import (
	"log"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type field2 struct {
	app.Compo
	Label        string
	Value        string
	Help         string
	HelpClass    string
	ReportChange func(v string) (string, error)
}

func newField2(v string, l string) *field2 {
	return &field2{
		Label: l,
		Value: v,
		Help:  v,
	}
}

func (f *field2) WithReportChange(fn func(v string) (string, error)) *field2 {
	f.ReportChange = fn
	return f
}

func (f *field2) Render() app.UI {
	return app.Div().Body(
		app.Label().Text(f.Label).Style("color", "blue"),
		app.Input().Type("text").Value(f.Value).OnInput(f.valueChanged),
		app.P().Style("color", f.HelpClass).Body(
			app.Text("Help: "),
			app.Text(f.Help),
		),
	)
}

func (f *field2) valueChanged(ctx app.Context, e app.Event) {
	f.Value = ctx.JSSrc.Get("value").String()
	if f.ReportChange != nil {
		v, err := f.ReportChange(f.Value)
		log.Printf("ReportChange field2 return: %s", v)
		if err == nil {
			f.HelpClass = "green"
		} else {
			f.HelpClass = "red"
		}
		f.Help = v
	}
}
