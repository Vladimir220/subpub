package subpub_lib

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestSimpleSubpub(t *testing.T) {
	subpub := NewSubPub()
	carOutput := make(chan string)
	defer close(carOutput)

	t.Log("Creating new subscriptions")
	carSubscription, err := subpub.Subscribe("car", func(msg any) { carOutput <- fmt.Sprint(msg) })
	if err != nil {
		t.Fatal("Subscribe error")
	}

	teaOutput := make(chan string)
	defer close(teaOutput)
	_, err = subpub.Subscribe("tea", func(msg any) { teaOutput <- fmt.Sprint(msg) })
	if err != nil {
		t.Fatal("Subscribe error")
	}

	// To make sure that all listeners go into standby mode
	time.Sleep(time.Second * 2)

	t.Log("Publishing data to channels")
	subpub.Publish("car", "Hello everynyan!")
	subpub.Publish("tea", "How are you?")
	subpub.Publish("tea", "I'm fine, thank you.")
	subpub.Publish("car", 1234567)

	t.Log("Receiving subscription data from chanel 1...")
	t.Log("{")
	for range 2 {
		select {
		case buf := <-carOutput:
			t.Log(buf)
		case <-time.After(time.Second * 1):
			t.Fatal("Listening error")
		}
	}
	t.Log("}")

	t.Log("Attempt to retrieve non-existent data from chanel 1")
	select {
	case buf := <-carOutput:
		t.Fatal("Unexpected data:", buf)
	case <-time.After(time.Second * 1):
	}

	t.Log("Receiving subscription data from chanel 2...")
	t.Log("{")
	for range 2 {
		select {
		case buf := <-teaOutput:
			t.Log(buf)
		case <-time.After(time.Second * 1):
			t.Fatal("Listening error")
		}
	}
	t.Log("}")

	t.Log("Attempt to retrieve non-existent data from chanel 2")
	select {
	case buf := <-teaOutput:
		t.Fatal("Unexpected data:", buf)
	case <-time.After(time.Second * 1):
	}

	t.Log("Unsubscribe from a channel")
	carSubscription.Unsubscribe()

	subpub.Publish("car", 321)
	select {
	case buf, exst := <-carOutput:
		if exst {
			t.Fatal("Unexpected data:", buf)
		}
	case <-time.After(time.Second * 1):
	}

	t.Log("Closing a subpub without canceling the context")
	ctx, cancel := context.WithCancel(context.Background())
	err = subpub.Close(ctx)
	if err == nil {
		t.Fatal("The subpub closed with the context alive")
	}

	t.Log("Closing a subpub with canceling the context")
	cancel()
	err = subpub.Close(ctx)
	if err != nil {
		t.Fatal("The subpub not closed")
	}

	t.Log("Checking if a subscription exists after closing a subpub")
	subpub.Publish("tea", 321)
	select {
	case buf, exst := <-teaOutput:
		if exst {
			t.Fatal("Unexpected data:", buf)
		}
	case <-time.After(time.Second * 1):
	}

}
