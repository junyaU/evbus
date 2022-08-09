package evbus

import (
	"testing"
)

func TestNew(t *testing.T) {
	if New() == nil {
		t.Fail()
	}
}

func TestSubscribe(t *testing.T) {
	bus := New()

	tests := []struct {
		v       interface{}
		wantErr bool
	}{
		{v: 0, wantErr: true},
		{v: "hello", wantErr: true},
		{v: func() {}, wantErr: false},
	}

	for _, test := range tests {
		err := bus.Subscribe("topic", test.v)

		if test.wantErr {
			if err == nil {
				t.Fail()
			}
		} else {
			if err != nil {
				t.Fail()
			}
		}
	}

	fn := func() {}
	bus.Subscribe("topic1", fn)
	if err := bus.Subscribe("topic1", fn); err == nil {
		t.Fail()
	}
}

func TestPublish(t *testing.T) {
	bus := New()

	v := []int{1, 2, 3, 4}

	bus.Subscribe("topic1", func(v *[]int) {
		*v = append(*v, 5)
	})

	bus.Publish("topic1", &v)

	if err := bus.Publish("topic2", "v"); err == nil {
		t.Fail()
	}

	bus.Wait()

	if len(v) != 5 {
		t.Fail()
	}
}

func TestUnSubscribe(t *testing.T) {
	bus := New()

	fn := func() {}

	bus.Subscribe("topic1", fn)

	if err := bus.UnSubscribe("topic1", fn); err != nil {
		t.Fail()
	}

	if err := bus.UnSubscribe("topic2", fn); err == nil {
		t.Fail()
	}

	if err := bus.UnSubscribe("topic1", func() {}); err == nil {
		t.Fail()
	}
}
