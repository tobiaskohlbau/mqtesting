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

package store

import (
	"errors"
)

var (
	ErrNotFound = errors.New("store: item(s) not found")
	ErrConflict = errors.New("store: item conflict")
)

type Service interface {
	Messages() MessageService
	Users() UserService
}

type MessageService interface {
	New(topic, payload string, messageID uint16) (int64, error)
	ListByTopic(topic string) ([]*Message, error)
	Get(id int64) (*Message, error)
	Delete(id int64) error
}

type UserService interface {
	New(provider, id string) (int64, error)
	Get(id int64) (*User, error)
	GetByAuth(provider, id string) (*User, error)
	SetName(id int64, name string) error
}
