package subpub_lib

import (
	"context"
	"errors"
	"sync"
	"time"
)

type SimpleSubpub struct {
	msgs          map[string][]any
	notifications map[string]*sync.Cond
	mu            *sync.RWMutex
	stopSignal    chan struct{}
}

func CreateSimpleSubpub() *SimpleSubpub {
	return &SimpleSubpub{msgs: make(map[string][]any, 10), notifications: make(map[string]*sync.Cond, 10), stopSignal: make(chan struct{}), mu: &sync.RWMutex{}}
}

func (sp *SimpleSubpub) checkSubject(subject string) {
	sp.mu.Lock()
	if _, exst := sp.msgs[subject]; !exst {
		sp.msgs[subject] = make([]any, 0, 5)
	}
	if _, exst := sp.notifications[subject]; !exst {
		sp.notifications[subject] = &sync.Cond{L: sp.mu}
	}
	sp.mu.Unlock()
}

func (sp *SimpleSubpub) Publish(subject string, msg any) error {
	sp.checkSubject(subject)

	sp.mu.Lock()
	sp.msgs[subject] = append(sp.msgs[subject], msg)
	sp.mu.Unlock()

	sp.notifications[subject].Broadcast()

	return nil
}

func (sp *SimpleSubpub) listener(subject string, cb MessageHandler, unsubscribeSignal chan struct{}, newDataSignal *sync.Cond) {
	sp.mu.RLock()
	msgId := len(sp.msgs[subject])
	sp.mu.RUnlock()

	for {
		newDataSignal.L.Lock()
		newDataSignal.Wait()
		newDataSignal.L.Unlock()

		select {
		case <-unsubscribeSignal:
			return
		case <-sp.stopSignal:
			close(unsubscribeSignal)
			return
		default:
		}

		sp.mu.RLock()
		for msgId < len(sp.msgs[subject]) {
			cb(sp.msgs[subject][msgId])
			msgId++
		}
		sp.mu.RUnlock()
	}
}

func (sp *SimpleSubpub) Subscribe(subject string, cb MessageHandler) (Subscription, error) {
	sp.checkSubject(subject)
	unsubscribeSignal := make(chan struct{})
	newDataSignal := sp.notifications[subject]

	go sp.listener(subject, cb, unsubscribeSignal, newDataSignal)

	subscription := &SimpleSubscription{unsubscribeSignal: unsubscribeSignal}
	return subscription, nil
}

func (sp SimpleSubpub) Close(ctx context.Context) error {
	select {
	case <-ctx.Done():
		close(sp.stopSignal)
		sp.mu.Lock()
		for _, s := range sp.notifications {
			s.Broadcast()
		}
		sp.mu.Unlock()
	case <-time.After(5 * time.Second):
		return errors.New("context is running, try later")
	}
	return nil
}
