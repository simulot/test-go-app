package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type message struct {
	id   int
	text string
}

type messages struct {
	l        []message
	onChange func()
}

func newMessages() *messages {
	ms := messages{}
	return &ms
}
func (ms messages) notify() {
	if ms.onChange != nil {
		ms.onChange()
	}
}
func (ms *messages) add(m message) {
	ms.l = append(ms.l, m)
	ms.notify()
}

func (ms *messages) dismiss(m message) {
	for i := range ms.l {
		if ms.l[i].id == m.id {
			copy(ms.l[i:], ms.l[i+1:])    // Shift ms.l[i+1:] left one index.
			ms.l[len(ms.l)-1] = message{} // Erase last element (write zero value).
			ms.l = ms.l[:len(ms.l)-1]     // Truncate slice.
			ms.notify()
			return
		}
	}
}

type page struct {
	app.Compo
}

func (p *page) Render() app.UI {
	return app.Table().Body(
		app.Tr().Body(
			app.Td().Style("valign", "baseline").Body(newFlat()),
			app.Td().Style("valign", "baseline").Body(newComposed()),
		),
	)

}

type flat struct {
	app.Compo
	messages *messages
	id       int
}

func newFlat() *flat {
	return &flat{
		messages: newMessages(),
	}
}

func (f *flat) OnMount(ctx app.Context) {
	f.messages.onChange = func() {
		ctx.Dispatch(func(ctx app.Context) {
			// nop
		})
	}
}

func (f *flat) Render() app.UI {
	return app.Div().
		Body(
			app.H1().
				Text("Notifications test flat"),
			app.P().Body(
				app.Button().Text("Send a notification").OnClick(f.send),
			),
			app.Div().ID("NOTIFICATIONS").Body(
				app.Range(f.messages.l).
					Slice(func(i int) app.UI {
						m := f.messages.l[i]
						return app.P().ID(fmt.Sprintf("P#%d", m.id)).
							Body(
								app.Button().
									ID(strconv.Itoa(m.id)).
									Body(
										app.Text(m.text),
									).
									OnClick(
										func(app.Context, app.Event) {
											f.messages.dismiss(m)
										}, m.id),
							)
					}),
			),
		)
}

func (f *flat) send(ctx app.Context, e app.Event) {
	f.id++

	s := fmt.Sprintf("Flat #%d", f.id)
	log.Printf("send message %s", s)
	f.messages.add(message{id: f.id, text: s})
}

type composed struct {
	app.Compo
	messages *messages
	id       int
}

func newComposed() *composed {
	return &composed{
		messages: newMessages(),
	}
}

func (c *composed) send(ctx app.Context, e app.Event) {
	c.id++

	s := fmt.Sprintf("Composed #%d", c.id)
	log.Printf("send message %s", s)
	c.messages.add(message{id: c.id, text: s})
}

func (c *composed) Render() app.UI {
	return app.Div().ID("COMPOSED").
		Body(
			app.H1().
				Text("Notifications test composed"),
			app.P().Body(
				app.Button().Text("Send a notification").OnClick(c.send),
			),
			newList(c.messages),
		)
}

type list struct {
	app.Compo
	messages *messages
}

func newList(messages *messages) *list {
	return &list{
		messages: messages,
	}
}

func (l *list) OnMount(ctx app.Context) {
	l.messages.onChange = func() {
		ctx.Dispatch(func(ctx app.Context) {
			// nop
		})
	}
}

func (l *list) Render() app.UI {
	return app.Div().ID("NOTIFICATIONS").Body(
		app.Range(l.messages.l).
			Slice(func(i int) app.UI {
				m := l.messages.l[i]
				return newNotification(m, func() { l.messages.dismiss(m) })
			}),
	)
}

type notification struct {
	app.Compo
	m         message
	dismissFn func()
}

func newNotification(message message, dismissFn func()) *notification {
	return &notification{
		m:         message,
		dismissFn: dismissFn,
	}
}

func (n *notification) Render() app.UI {
	return app.P().ID(fmt.Sprintf("P#%d", n.m.id)).
		Body(
			app.Button().
				ID(strconv.Itoa(n.m.id)).
				Body(
					app.Text(n.m.text),
				).
				OnClick(
					func(app.Context, app.Event) {
						n.dismissFn()
					}, n.m.id),
		)
}
