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

import "kohlbau.de/x/mqtesting/store"

type Store struct {
	messageService *messageService
	userService    *userService
}

func (s *Store) Messages() store.MessageService {
	return s.messageService
}

func (s *Store) Users() store.UserService {
	return s.userService
}

func New() *Store {
	return &Store{
		messageService: &messageService{
			msgs: make(map[int64]*store.Message),
		},
		userService: &userService{
			users: make(map[int64]*store.User),
		},
	}
}
