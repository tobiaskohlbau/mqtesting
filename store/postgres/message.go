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

type messageStore struct {
	db *sql.DB
}

func (s *messageStore) New(topic, payload string, messageID uint16) (int64, error) {
	var id int64

	err := s.db.QueryRow(
		`insert into messages(topic, payload, message_id, created_at) values($1, $2, $3, $4) returning id`,
		topic,
		payload,
		messageID,
		time.Now(),
	).Scan(&id)

	return id, err
}

const selectFromMessages = `select * from messages`

func (s *messageStore) scanMessage(scanner scanner) (*store.Message, error) {
	m := new(store.Message)

	err := scanner.Scan(&m.ID, &m.Duplicate, &m.Qos, &m.Retained, &m.Topic, &m.MessageID, &m.Payload, &m.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (s *messageStore) ListByTopic(topic string) ([]*store.Message, error) {
	var (
		msgs []*store.Message
		rows *sql.Rows
		err  error
	)

	if topic == "" {
		rows, err = s.db.Query(selectFromMessages)
	} else {
		rows, err = s.db.Query(selectFromMessages+` where topic=$1`, topic)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		m, err := s.scanMessage(rows)
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return msgs, nil
}

func (s *messageStore) Get(id int64) (*store.Message, error) {
	row := s.db.QueryRow(selectFromMessages+` where id=$1`, id)
	return s.scanMessage(row)
}

func (s *messageStore) Delete(id int64) error {
	_, err := s.db.Exec(`delete from messages where id=$1`, id)
	return err
}
