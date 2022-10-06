package resolvers

import (
	"context"

	"github.com/asaskevich/EventBus"
)

var bus = EventBus.New()

func subscribeUntilDone(ctx context.Context, topic string, eventHandler interface{}) error {
	// Execute eventHandler for every message on topic.
	err := bus.Subscribe(topic, eventHandler)
	if err != nil {
		return err
	}

	// Launch subroutine that will block until context is done (which is the
	// end of the GraphQL subscription), after which we unsubscribe.
	go func() {
		select {
		case <-ctx.Done():
			bus.Unsubscribe(topic, eventHandler)
			return
		}
	}()

	return nil
}
