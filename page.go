package main

import (
	"errors"
	"log"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type page struct {
	app.Compo

	Field1 string
	Field2 string
}

func (p *page) OnMount(ctx app.Context) {
	p.Field1 = "Initial value field1"
	p.Field2 = "Initial value field2"

}

func (p *page) Render() app.UI {
	return app.Div().Body(
		newField1(p.Field1, "Field 1").WithReportChange(func(v string) (string, error) {
			log.Printf("Reported Field 1 value: %s", v)
			p.Field1 = v
			if len(v)%3 == 0 {
				return "You made an error", errors.New("error")
			}
			return "All good", nil
		}),

		newField2(p.Field2, "Field 2").WithReportChange(func(v string) (string, error) {
			log.Printf("Reported Field 2 value: %s", v)
			p.Field2 = v
			if len(v)%3 == 0 {
				return "You made an error", errors.New("error")
			}
			return "All good", nil
		}),
		app.Hr(),
		app.Text("Reported Field 1 "),
		app.Text(p.Field1),

		app.Br(),
		app.Text("Reported Field 2 "),
		app.Text(p.Field2),
	)
}
