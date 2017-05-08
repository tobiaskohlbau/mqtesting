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

package mq

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/tobiaskohlbau/mqtesting/store"
)

type Mq struct {
	store store.Store
}

func (m *Mq) onMessage(c mqtt.Client, msg mqtt.Message) {
	_, err := m.store.Messages().New(msg.Topic(), string(msg.Payload()[:]), msg.MessageID())
	if err != nil {
		log.Fatal(err)
	}
}

func New(address string, store store.Store) *Mq {
	mq := &Mq{
		store: store,
	}

	opts := &mqtt.ClientOptions{
		ClientID:     "mqtesting-sub",
		CleanSession: true,
		OnConnect: func(c mqtt.Client) {
			if token := c.Subscribe("#", 0, mq.onMessage); token.Wait() && token.Error() != nil {
				log.Fatal(token.Error())
			}
		},
	}
	opts.AddBroker(fmt.Sprintf("tcp://%s", address))
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	return mq
}
