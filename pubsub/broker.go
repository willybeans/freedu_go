package pubsub

import (
	"sync"

	"github.com/willybeans/freedu_go/logger"
)

type Subscribers map[string]*Subscriber

type Broker struct {
	subscribers Subscribers            // map of subscribers id:Subscriber
	topics      map[string]Subscribers // map of topic to subscribers
	mut         sync.RWMutex           // mutex lock
}

func NewBroker() *Broker {
	// returns new broker object
	return &Broker{
		subscribers: Subscribers{},
		topics:      map[string]Subscribers{},
	}
}

func (b *Broker) AddSubscriber(user_uuid string) *Subscriber {
	b.mut.Lock()
	defer b.mut.Unlock()
	// pass id
	id, s := CreateNewSubscriber(user_uuid)
	b.subscribers[id] = s
	return s
}

func (b *Broker) RemoveSubscriber(s *Subscriber) {
	for topic := range s.topics {
		b.Unsubscribe(s, topic)
	}
	b.mut.Lock()
	delete(b.subscribers, s.id)
	b.mut.Unlock()
	s.Destruct()
}

func (b *Broker) Broadcast(msg string, topics []string) {
	// broadcast message to all topics.
	for _, topic := range topics {
		for _, s := range b.topics[topic] {
			m := NewMessage(msg, topic)
			go (func(s *Subscriber) {
				s.Signal(m)
			})(s)
		}
	}
}

func (b *Broker) GetSubscribers(topic string) int {
	b.mut.RLock() // whats the diff with .Lock()? whats recursive read locking?
	defer b.mut.RUnlock()
	return len(b.topics[topic])
}

func (b *Broker) Subscribe(s *Subscriber, topic string) {
	l := logger.Get()
	b.mut.Lock()
	defer b.mut.Unlock()

	if b.topics[topic] == nil {
		b.topics[topic] = Subscribers{} // i think this is an empty map?
	}
	s.AddTopic(topic)
	b.topics[topic][s.id] = s
	l.Info().Msgf("%s Subscribed for topic: %s\n", s.id, topic)
}

func (b *Broker) Unsubscribe(s *Subscriber, topic string) {
	l := logger.Get()
	b.mut.RLock()
	defer b.mut.RUnlock()

	delete(b.topics[topic], s.id)
	s.RemoveTopic(topic)
	l.Info().Msgf("%s Unsubscribed for topic: %s\n", s.id, topic)
}

func (b *Broker) Publish(topic string, msg string) {
	// publish the message to given topic.
	b.mut.RLock()
	bTopics := b.topics[topic]
	b.mut.RUnlock()
	for _, s := range bTopics {
		m := NewMessage(msg, topic)
		if !s.active {
			return
		}
		go (func(s *Subscriber) {
			s.Signal(m)
		})(s)
	}
}
