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

package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"kohlbau.de/x/mqtesting/jwt"
	"kohlbau.de/x/mqtesting/oauth"
	"kohlbau.de/x/mqtesting/store"

	"github.com/goware/cors"
	"github.com/pressly/chi"
)

type Service interface {
	http.Handler
}

type service struct {
	rtr   chi.Router
	store store.Service
	jwt   jwt.Service
	oauth oauth.Service
}

func New(opts ...Option) Service {
	s := &service{
		rtr: chi.NewRouter(),
		jwt: jwt.New(),
	}

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	s.rtr.Use(cors.Handler)

	for _, opt := range opts {
		opt(s)
	}

	s.rtr.Use(s.handleNoStore)
	s.rtr.Get("/messages", s.handleGetMessages)
	s.rtr.Get("/messages/:id", s.handleGetMessage)
	s.rtr.Delete("/messages/:id", s.handleDeleteMessage)

	s.rtr.Get("/user", s.handleUser)

	s.rtr.Get("/provider", s.handleProvider)

	return s
}

type Option func(*service)

func WithOAuth(oauth oauth.Service) Option {
	return func(s *service) {
		s.oauth = oauth
	}
}

func WithJWT(jwt jwt.Service) Option {
	return func(s *service) {
		s.jwt = jwt
	}
}

func WithStore(store store.Service) Option {
	return func(s *service) {
		s.store = store
	}
}

func (s *service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.rtr.ServeHTTP(w, r)
}

func (s *service) urlParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

func (s *service) user(r *http.Request) *store.User {
	ah := r.Header.Get("Authorization")
	if len(ah) < 7 || strings.ToUpper(ah[:6]) != "BEARER" {
		return nil
	}
	t := ah[7:]

	uid, _, err := s.jwt.Verify(t)
	if err != nil {
		return nil
	}

	usr, err := s.store.Users().Get(uid)
	if err != nil {
		return nil
	}

	return usr
}

func (s *service) render(w http.ResponseWriter, st int, d interface{}) {
	b, err := json.Marshal(d)
	if err != nil {
		http.Error(w, "failed to marshal data", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(st)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func (s *service) handleProvider(w http.ResponseWriter, r *http.Request) {
	pvds := []string{}
	if s.oauth != nil {
		pvds = s.oauth.Provider()
	}
	b, err := json.Marshal(pvds)
	if err != nil {
		http.Error(w, "failed to marshal provider", http.StatusInternalServerError)
	}
	w.Write(b)
}

func (s *service) handleNoStore(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.store == nil {
			http.Error(w, "no storage specified", http.StatusInternalServerError)
			return
		}
		next.ServeHTTP(w, r)
	})
}
