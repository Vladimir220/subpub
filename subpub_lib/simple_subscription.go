package subpub_lib

import "sync"

type SimpleSubscription struct {
	unsubscribeSignal chan struct{}
	subSignal         *sync.Cond
}

func (subscription *SimpleSubscription) Unsubscribe() {
	close(subscription.unsubscribeSignal)
	subscription.subSignal.Broadcast()
}
