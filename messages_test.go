package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

func TestMessages(t *testing.T) {
	t.Run("add message", func(t *testing.T) {
		notified := false
		notifyMe := func() {
			notified = true
		}

		mess := newMessages()
		mess.onChange = notifyMe

		want := message{id: 42, text: "the answere"}
		mess.add(want)

		if !notified {
			t.Errorf("Expected to be notified, but not")
		}
	})
	t.Run("add 3 messages, delete 1st", func(t *testing.T) {
		notified := false
		notifyMe := func() {
			notified = true
		}

		list := []message{
			{id: 1, text: "1"},
			{id: 2, text: "2"},
			{id: 3, text: "3"},
		}

		mess := newMessages()
		mess.onChange = notifyMe

		for _, m := range list {
			mess.add(m)
		}

		if !notified {
			t.Errorf("Expected to be notified, but not")
		}

		if len(mess.l) != 3 {
			t.Errorf("Expected 3 messages, got %d", len(mess.l))
		}

		notified = false
		mess.dismiss(list[0])

		if !notified {
			t.Errorf("Expected to be notified after dismess, but not")
		}

		if len(mess.l) != 2 {
			t.Errorf("Expected 2 messages, got %d", len(mess.l))
		}

		want := []message{
			// {id: 1, text: "1"},
			{id: 2, text: "2"},
			{id: 3, text: "3"},
		}

		if !reflect.DeepEqual(want, mess.l) {
			t.Errorf("Expected %#v, got %#v", want, mess.l)
		}
	})

}

type testcase struct {
	app.TestUIDescriptor
	expectedErr bool
}

func testMatches(t *testing.T, got app.UI, expected []testcase) error {
	t.Helper()
	for _, e := range expected {
		err := app.TestMatch(got, e.TestUIDescriptor)
		if e.expectedErr && (err == nil) {
			return fmt.Errorf("Expected error but got no error")
		}
		if !e.expectedErr && (err != nil) {
			return err
		}
	}
	return nil
}

func TestFlat(t *testing.T) {
	list := []message{
		{id: 0, text: "1"},
		{id: 1, text: "2"},
		{id: 2, text: "3"},
	}

	t.Run("Check that all messages are displayed", func(t *testing.T) {
		c := newFlat()
		disp := app.NewClientTester(c)
		defer disp.Close()
		disp.Consume()

		onChangeCalled := false

		c.messages.onChange = func(fn func()) func() {
			return func() {
				onChangeCalled = true
				fn()
			}
		}(c.messages.onChange)

		err := app.TestMatch(c,
			app.TestUIDescriptor{
				Path:     app.TestPath(0, 2),
				Expected: app.Div().ID("NOTIFICATIONS"),
			})
		if err != nil {
			t.Errorf("Expecting  Div#NOTIFICATIONS: got %s", err)
			return
		}
		for i := 0; i < len(list); i++ {
			c.messages.add(list[i])
		}
		disp.Consume()

		if !onChangeCalled {
			t.Errorf("Expecting message.onChange called")
		}

		err = testMatches(t, c, []testcase{
			{
				TestUIDescriptor: app.TestUIDescriptor{
					Path:     app.TestPath(0, 2, 0),
					Expected: app.P().ID("P#0"),
				},
				expectedErr: false,
			},
			{
				TestUIDescriptor: app.TestUIDescriptor{
					Path:     app.TestPath(0, 2, 1),
					Expected: app.P().ID("P#1"),
				},
				expectedErr: false,
			},
			{
				TestUIDescriptor: app.TestUIDescriptor{
					Path:     app.TestPath(0, 2, 2),
					Expected: app.P().ID("P#2"),
				},
				expectedErr: false,
			},
		})
		if err != nil {
			t.Errorf("UI not as expected: %s", err)
			return
		}
	})
	t.Run("Test dimiss last message", func(t *testing.T) {
		c := newFlat()
		disp := app.NewClientTester(c)
		defer disp.Close()
		disp.Consume()
		for i := 0; i < len(list); i++ {
			c.messages.add(list[i])
		}
		disp.Consume()

		c.messages.dismiss(list[len(list)-1])
		disp.Consume()

		err := testMatches(t, c,
			[]testcase{
				{
					TestUIDescriptor: app.TestUIDescriptor{
						Path:     app.TestPath(0, 2, 0),
						Expected: app.P().ID("P#0"),
					},
					expectedErr: false,
				},
				{
					TestUIDescriptor: app.TestUIDescriptor{
						Path:     app.TestPath(0, 2, 1),
						Expected: app.P().ID("P#1"),
					},
					expectedErr: false,
				},

				{
					TestUIDescriptor: app.TestUIDescriptor{
						Path:     app.TestPath(0, 2, 2),
						Expected: app.P().ID("P#2"),
					},
					expectedErr: true,
				},
			})
		if err != nil {
			t.Errorf("UI not as expected: %s", err)
			return
		}
	})

	t.Run("Test dimiss first message", func(t *testing.T) {
		c := newFlat()
		disp := app.NewClientTester(c)
		defer disp.Close()
		disp.Consume()
		for i := 0; i < len(list); i++ {
			c.messages.add(list[i])
		}
		disp.Consume()

		c.messages.dismiss(list[0])
		disp.Consume()
		err := testMatches(t, c,
			[]testcase{
				{
					TestUIDescriptor: app.TestUIDescriptor{
						Path:     app.TestPath(0, 2, 0),
						Expected: app.P().ID("P#1"),
					},
					expectedErr: false,
				},
				{
					TestUIDescriptor: app.TestUIDescriptor{
						Path:     app.TestPath(0, 2, 1),
						Expected: app.P().ID("P#2"),
					},
					expectedErr: false,
				},

				{
					TestUIDescriptor: app.TestUIDescriptor{
						Path:     app.TestPath(0, 2, 2),
						Expected: app.P().ID("P#3"),
					},
					expectedErr: true,
				},
			})
		if err != nil {
			t.Errorf("UI not as expected: %s", err)
			return
		}
	})

}

func TestComposed(t *testing.T) {
	list := []message{
		{id: 0, text: "1"},
		{id: 1, text: "2"},
		{id: 2, text: "3"},
	}

	t.Run("Check that all messages are displayed", func(t *testing.T) {
		c := newComposed()
		disp := app.NewClientTester(c)
		defer disp.Close()
		disp.Consume()

		onChangeCalled := false

		c.messages.onChange = func(fn func()) func() {
			return func() {
				onChangeCalled = true
				fn()
			}
		}(c.messages.onChange)

		err := app.TestMatch(c,
			app.TestUIDescriptor{
				Path:     app.TestPath(0, 2, 0),
				Expected: app.Div().ID("NOTIFICATIONS"),
			})
		if err != nil {
			t.Errorf("Expecting  Div#NOTIFICATIONS: got %s", err)
			return
		}
		for i := 0; i < len(list); i++ {
			c.messages.add(list[i])
		}
		disp.Consume()

		if !onChangeCalled {
			t.Errorf("Expecting message.onChange called")
		}
		// Path:
		// 0 -> div#COMPOSED
		// 0,0 --> <H1>Notifications test composed</H1>
		// 0,1 --> <p>
		// 0,2 --> Component list
		// 0,2,0 --> div#NOTIFICATIONS
		// 0,2,0,0 --> component message
		// 0,2,0,0,0 --> P#0
		// 0,2,0,1,0 --> P#1

		err = testMatches(t, c, []testcase{
			{
				TestUIDescriptor: app.TestUIDescriptor{
					Path:     app.TestPath(0, 2, 0, 0, 0),
					Expected: app.P().ID("P#0"),
				},
				expectedErr: false,
			},
			{
				TestUIDescriptor: app.TestUIDescriptor{
					Path:     app.TestPath(0, 2, 0, 1, 0),
					Expected: app.P().ID("P#1"),
				},
				expectedErr: false,
			},
			{
				TestUIDescriptor: app.TestUIDescriptor{
					Path:     app.TestPath(0, 2, 0, 2, 0),
					Expected: app.P().ID("P#2"),
				},
				expectedErr: false,
			},
		})
		if err != nil {
			t.Errorf("UI not as expected: %s", err)
			return
		}
	})
	t.Run("Test dimiss last message", func(t *testing.T) {
		c := newComposed()
		disp := app.NewClientTester(c)
		defer disp.Close()
		disp.Consume()
		for i := 0; i < len(list); i++ {
			c.messages.add(list[i])
		}
		disp.Consume()

		c.messages.dismiss(list[len(list)-1])
		disp.Consume()

		err := testMatches(t, c,
			[]testcase{
				{
					TestUIDescriptor: app.TestUIDescriptor{
						Path:     app.TestPath(0, 2, 0, 0, 0),
						Expected: app.P().ID("P#0"),
					},
					expectedErr: false,
				},
				{
					TestUIDescriptor: app.TestUIDescriptor{
						Path:     app.TestPath(0, 2, 0, 1, 0),
						Expected: app.P().ID("P#1"),
					},
					expectedErr: false,
				},

				{
					TestUIDescriptor: app.TestUIDescriptor{
						Path:     app.TestPath(0, 2, 0, 2, 0),
						Expected: app.P().ID("P#2"),
					},
					expectedErr: true,
				},
			})
		if err != nil {
			t.Errorf("UI not as expected: %s", err)
			return
		}
	})

	t.Run("Test dimiss first message", func(t *testing.T) {
		c := newComposed()
		disp := app.NewClientTester(c)
		defer disp.Close()
		disp.Consume()
		for i := 0; i < len(list); i++ {
			c.messages.add(list[i])
		}
		disp.Consume()

		c.messages.dismiss(list[0])
		disp.Consume()

		err := testMatches(t, c,
			[]testcase{
				{
					TestUIDescriptor: app.TestUIDescriptor{
						Path:     app.TestPath(0, 2, 0, 0, 0),
						Expected: app.P().ID("P#1"),
					},
					expectedErr: false,
				},
				{
					TestUIDescriptor: app.TestUIDescriptor{
						Path:     app.TestPath(0, 2, 0, 1, 0),
						Expected: app.P().ID("P#2"),
					},
					expectedErr: false,
				},

				{
					TestUIDescriptor: app.TestUIDescriptor{
						Path:     app.TestPath(0, 2, 0, 2, 0),
						Expected: app.P().ID("P#2"),
					},
					expectedErr: true,
				},
			})
		if err != nil {
			t.Errorf("UI not as expected: %s", err)
			return
		}
	})

}
