package notifier

type FooPublisher struct {
	Subscribers []IFooSubscriber
}

// Subscribe to publisher
func (p *FooPublisher) Subscribe(subscriber IFooSubscriber) {
	p.Subscribers = append(p.Subscribers, subscriber)
}

// Unsubscribe to publisher
func (p *FooPublisher) Unsubscribe(subscriber IFooSubscriber) {
	for i := len(p.Subscribers) - 1; i >= 0; i-- {
		if p.Subscribers[i] == subscriber {
			p.Subscribers = append(p.Subscribers[:i], p.Subscribers[i+1:]...)
		}
	}
}

func NewPublisher() *FooPublisher {
	return &FooPublisher{Subscribers: make([]IFooSubscriber, 0)}
}
