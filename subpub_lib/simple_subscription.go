package subpub_lib

type SimpleSubscription struct {
	unsubscribeSignal chan struct{}
}

func (subscription *SimpleSubscription) Unsubscribe() {
	close(subscription.unsubscribeSignal)
}
