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
	"time"

	"github.com/tobiaskohlbau/mqtesting/store"
)

type userService struct {
	db *sql.DB
}

const selectFromUsers = `
	select
		id,
		coalesce(name, '') as name,
		created_at,
		auth_provider,
		auth_id
    from users
`

func (s *userService) scanUser(scanner scanner) (*store.User, error) {
	u := &store.User{}
	err := scanner.Scan(
		&u.ID,
		&u.Name,
		&u.CreatedAt,
		&u.AuthProvider,
		&u.AuthID,
	)
	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *userService) New(provider, id string) (int64, error) {
	var uid int64
	err := s.db.QueryRow(`insert into users(created_at, auth_provider, auth_id) values($1, $2, $3) returning id`,
		time.Now(), provider, id,
	).Scan(&id)
	return uid, err
}

func (s *userService) Get(id int64) (*store.User, error) {
	row := s.db.QueryRow(selectFromUsers+`where id=$1`, id)
	return s.scanUser(row)
}

func (s *userService) GetByAuth(provider, id string) (*store.User, error) {
	row := s.db.QueryRow(selectFromUsers+` where auth_provider=$1 and auth_id=$2`, provider, id)
	return s.scanUser(row)
}

func (s *userService) SetName(id int64, name string) error {
	_, err := s.db.Exec("update users set name=$1 where id=$2", name, id)
	return err
}
