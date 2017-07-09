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

package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"kohlbau.de/x/mqtesting/store"
)

type Store struct {
	db             *sql.DB
	messageService *messageService
	userService    *userService
}

func (s *Store) Messages() store.MessageService {
	return s.messageService
}

func (s *Store) Users() store.UserService {
	return s.userService
}

func Connect(address, username, password, database string) (*Store, error) {
	cstr := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		username, password, address, database,
	)

	db, err := sql.Open("postgres", cstr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	s := &Store{
		db:             db,
		messageService: &messageService{db: db},
		userService:    &userService{db: db},
	}

	err = s.Migrate()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Store) Migrate() error {
	for _, m := range migrate {
		if _, err := s.db.Exec(m); err != nil {
			return fmt.Errorf("error executing sql: %s: query: %q", err, m)
		}
	}
	return nil
}

func (s *Store) Drop() error {
	for _, d := range drop {
		if _, err := s.db.Exec(d); err != nil {
			return fmt.Errorf("error executing sql: %s: query: %q", err, d)
		}
	}
	return nil
}

func (s *Store) Reset() error {
	if err := s.Drop(); err != nil {
		return err
	}
	return s.Migrate()
}

type scanner interface {
	Scan(v ...interface{}) error
}
