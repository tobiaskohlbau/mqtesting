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

package memory

import (
	"sort"

	"time"

	"github.com/tobiaskohlbau/mqtesting/store"
)

type messageService struct {
	msgs map[int64]*store.Message
	id   int64
}

func (s *messageService) New(topic, payload string, messageID uint16) (int64, error) {
	defer func() { s.id++ }()
	s.msgs[s.id] = &store.Message{
		ID:        s.id,
		Topic:     topic,
		Payload:   payload,
		MessageID: messageID,
		CreatedAt: time.Now(),
	}
	return s.id, nil
}

func (s *messageService) ListByTopic(topic string) ([]*store.Message, error) {
	msgs := []*store.Message{}
	for _, msg := range s.msgs {
		if topic == "" || msg.Topic == topic {
			msgs = append(msgs, msg)
		}
	}
	sort.Slice(msgs, func(i, j int) bool { return msgs[i].ID < msgs[j].ID })
	return msgs, nil
}

func (s *messageService) Get(id int64) (*store.Message, error) {
	msg, ok := s.msgs[id]
	if !ok {
		return nil, store.ErrNotFound
	}
	return msg, nil
}

func (s *messageService) Delete(id int64) error {
	msg, err := s.Get(id)
	if err != nil {
		return err
	}
	s.msgs[msg.ID] = nil
	return nil
}
