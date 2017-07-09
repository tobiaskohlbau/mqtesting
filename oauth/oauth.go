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

package oauth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"kohlbau.de/x/mqtesting/jwt"
	"kohlbau.de/x/mqtesting/store"

	"golang.org/x/oauth2"

	"github.com/pressly/chi"
	uuid "github.com/satori/go.uuid"
)

const (
	stateCookie = "mqtesting_oauth_state"
	tokenCookie = "mqtesting_jwt_token"
	timeout     = 10 * time.Second
)

type Service interface {
	Provider() []string
	http.Handler
}

type service struct {
	r           chi.Router
	providers   map[string]*provider
	str         store.UserService
	jwt         jwt.Service
	baseURL     string
	frontendURL string
}

type Option func(*service)

func New(userService store.UserService, jwtService jwt.Service, baseURL string, opts ...Option) Service {
	s := &service{
		r:           chi.NewRouter(),
		providers:   make(map[string]*provider),
		str:         userService,
		jwt:         jwtService,
		baseURL:     baseURL,
		frontendURL: baseURL,
	}

	for _, opt := range opts {
		opt(s)
	}

	s.r.Get("/authenticate/:provider", s.handleAuthenticate)
	s.r.Get("/callback/:provider", s.handleCallback)

	return s
}

func WithProvider(name, id, secret string) Option {
	c, ok := providerConfigs[name]
	if !ok {
		log.Printf("oauth: provider not supported: %q", name)
		return func(*service) {}
	}
	return func(s *service) {
		s.providers[name] = &provider{
			cfg: &oauth2.Config{
				ClientID:     id,
				ClientSecret: secret,
				RedirectURL:  fmt.Sprintf("%s/callback/%s", s.baseURL, name),
				Endpoint:     c.endpoint,
				Scopes:       c.scopes,
			},
			user: c.user,
		}
	}
}

func WithFrontend(url string) Option {
	return func(s *service) {
		s.frontendURL = url
	}
}

func (s *service) Provider() []string {
	pvds := []string{}
	for name := range s.providers {
		pvds = append(pvds, name)
	}
	sort.Slice(pvds, func(i, j int) bool { return pvds[i] < pvds[j] })
	return pvds
}

func (s *service) handleAuthenticate(w http.ResponseWriter, r *http.Request) {
	pn := chi.URLParam(r, "provider")
	p, ok := s.providers[pn]
	if !ok {
		http.NotFound(w, r)
		return
	}

	id := uuid.NewV4().String()

	http.SetCookie(w, &http.Cookie{
		Name:     stateCookie,
		Value:    id,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   3600,
	})

	http.Redirect(w, r, p.cfg.AuthCodeURL(id), http.StatusFound)
}

func (s *service) handleCallback(w http.ResponseWriter, r *http.Request) {
	pn := chi.URLParam(r, "provider")
	p, ok := s.providers[pn]
	if !ok {
		http.NotFound(w, r)
		return
	}

	c, err := r.Cookie(stateCookie)
	if err != nil {
		s.Error(w, "failed to get oauth state cookie", err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     stateCookie,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   0,
	})

	state := c.Value
	if state == "" {
		s.Error(w, "empty oauth state cookie")
		return
	}

	qs := r.URL.Query().Get("state")
	if qs != state {
		s.Error(w, "bad state value")
		return
	}

	qc := r.URL.Query().Get("code")
	if qc == "" {
		s.Error(w, "empty code value")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeout)
	defer cancel()

	t, err := p.cfg.Exchange(ctx, qc)
	if err != nil {
		s.Error(w, "exchange failed: %v", err)
		return
	}
	if !t.Valid() {
		s.Error(w, "invalid token")
		return
	}

	id, name, err := p.user(p.cfg.Client(ctx, t))
	if err != nil {
		s.Error(w, "get provider user failed: %v", err)
		return
	}

	if id == "" {
		s.Error(w, "provider user id empty")
		return
	}

	var token string

	u, err := s.str.GetByAuth(pn, id)
	switch err {
	case nil:
		token, err = s.jwt.Create(u.ID)
		if err != nil {
			s.Error(w, "failed to create auth token: %v", err)
			return
		}
	case store.ErrNotFound:
		uid, err := s.str.New(pn, id)
		if err != nil {
			s.Error(w, "failed to create user: %v", err)
			return
		}

		token, err = s.jwt.Create(uid)
		if err != nil {
			s.Error(w, "failed to create auth token: %v", err)
			return
		}
	default:
		s.Error(w, "failed to get user by auth: %v", err)
		return
	}

	if err := s.str.SetName(u.ID, name); err != nil {
		http.Error(w, "oauth: failed to set username", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "<!doctype html><title>OAuth Successfull</title><script>var message = {\"status\": %d, \"jwt\": \"%s\"};window.opener.postMessage(message, '%s');window.close();</script>", 200, token, s.frontendURL)
}

func (s *service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.r.ServeHTTP(w, r)
}

func (s *service) Error(w http.ResponseWriter, msg string, f ...interface{}) {
	fmt.Fprintf(w, msg, f...)
	w.WriteHeader(http.StatusInternalServerError)
}
