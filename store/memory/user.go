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

type userService struct {
	users map[int64]*store.User
	id    int64
}

func (s *userService) New(provider, id string) (int64, error) {
	_, err := s.GetByAuth(provider, id)
	if err == store.ErrNotFound {
		s.id++
		s.users[s.id] = &store.User{
			ID:           s.id,
			AuthProvider: provider,
			AuthID:       id,
		}
		return s.id, nil
	}
	return -1, store.ErrConflict
}

func (s *userService) Get(id int64) (*store.User, error) {
	u, ok := s.users[id]
	if !ok {
		return nil, store.ErrNotFound
	}
	return u, nil
}

func (s *userService) GetByAuth(provider, id string) (*store.User, error) {
	for _, u := range s.users {
		if u.AuthProvider == provider && u.AuthID == id {
			return u, nil
		}
	}
	return nil, store.ErrNotFound
}

func (s *userService) SetName(id int64, name string) error {
	u, ok := s.users[id]
	if !ok {
		return store.ErrNotFound
	}
	u.Name = name
	return nil
}
