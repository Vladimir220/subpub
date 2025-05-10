package subpub_lib

import (
	"context"
	"sync"

	"github.com/redis/go-redis/v9"
)

type RedisSubpub struct {
	rdb        *redis.Client
	ctx        context.Context
	stopSignal chan struct{}
}

func CreateRedisSubpub() *RedisSubpub {
	return &RedisSubpub{stopSignal: make(chan struct{})}
}

func (sp *RedisSubpub) Publish(subject string, msg any) error {

	return nil
}

func (sp *RedisSubpub) listener(subject string, cb MessageHandler, unsubscribeSignal chan struct{}, newDataSignal *sync.Cond) {
	/*sp.mu.RLock()
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
	}*/
}

func (sp *RedisSubpub) Subscribe(subject string, cb MessageHandler) (Subscription, error) {
	unsubscribeSignal := make(chan struct{})

	//go sp.listener(subject, cb, unsubscribeSignal, newDataSignal)

	subscription := &SimpleSubscription{unsubscribeSignal: unsubscribeSignal}

	return subscription, nil
}

func (sp RedisSubpub) Close(ctx context.Context) error {
	/*select {
	case <-ctx.Done():
		close(sp.stopSignal)
		sp.mu.Lock()
		for _, s := range sp.notifications {
			s.Broadcast()
		}
		sp.mu.Unlock()
	case <-time.After(5 * time.Second):
		return errors.New("context is running, try later")
	}*/
	return nil
}
