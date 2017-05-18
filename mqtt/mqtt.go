//    Copyright 2017 Tobias Kohlbau
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package mqtt

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Client mqtt.Client
type Message mqtt.Message
type Handler func(client Client, msg Message)

type Service interface {
	Publish(topic string, qos byte, retained bool, payload interface{}) error
	Subscribe(topic string, handler Handler) error
}

type service struct {
	c mqtt.Client
}

func New(broker string, id string) (Service, error) {
	mq := &service{}

	opts := &mqtt.ClientOptions{
		ClientID:     id,
		CleanSession: true,
	}
	opts.AddBroker(fmt.Sprintf("tcp://%s", broker))
	mq.c = mqtt.NewClient(opts)

	if t := mq.c.Connect(); t.Wait() && t.Error() != nil {
		return nil, t.Error()
	}

	return mq, nil
}

func (s *service) Subscribe(topic string, handler Handler) error {
	h := func(c mqtt.Client, msg mqtt.Message) {
		handler(c, msg)
	}
	if t := s.c.Subscribe(topic, 0, h); t.Wait() && t.Error() != nil {
		return t.Error()
	}
	return nil
}

func (s *service) Publish(topic string, qos byte, retained bool, payload interface{}) error {
	if t := s.c.Publish(topic, qos, retained, payload); t.Wait() && t.Error() != nil {
		return t.Error()
	}
	return nil
}
