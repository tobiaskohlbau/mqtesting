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

package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Service interface {
	Create(id int64) (token string, err error)
	Verify(token string) (id int64, issuedAt time.Time, err error)
}

type Option func(s *service)

func WithSecret(secret []byte) Option {
	return func(s *service) {
		s.secret = secret
	}
}

func New(opts ...Option) Service {
	s := &service{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

type service struct {
	secret []byte
}

type claims struct {
	jwt.StandardClaims
	UserID *int64 `json:"_uid"`
}

func (s *service) Create(id int64) (string, error) {
	claims := claims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt: time.Now().Unix(),
		},
		UserID: &id,
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ts, err := t.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("jwt: token signing failed: %v", err)
	}

	return ts, nil
}

func (s *service) Verify(ts string) (int64, time.Time, error) {
	t, err := jwt.ParseWithClaims(
		ts,
		&claims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("jwt: unexpected signing method")
			}
			return s.secret, nil
		},
	)
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("jwt: ParseWithClaims failed: %v", err)
	}
	if !t.Valid {
		return 0, time.Time{}, errors.New("jwt: token not valid")
	}
	c, ok := t.Claims.(*claims)
	if !ok {
		return 0, time.Time{}, errors.New("jwt: failed to get token claims")
	}

	if c.UserID == nil {
		return 0, time.Time{}, errors.New("jwt: UserID claim is not valid")
	}
	if c.IssuedAt == 0 {
		return 0, time.Time{}, errors.New("jwt: IssuedAt claim is not valid")
	}

	return *c.UserID, time.Unix(c.IssuedAt, 0), nil
}
