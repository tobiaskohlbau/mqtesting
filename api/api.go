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

	"github.com/goware/cors"
	"github.com/pressly/chi"
	"github.com/tobiaskohlbau/mqtesting/mq"
	"github.com/tobiaskohlbau/mqtesting/store"
)

type api struct {
	rtr chi.Router
	str store.Store
	mq  *mq.Mq
}

func New(store store.Store, opts ...ApiOption) http.Handler {
	a := &api{
		rtr: chi.NewRouter(),
		str: store,
	}

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	a.rtr.Use(cors.Handler)

	for _, opt := range opts {
		opt(a)
	}

	a.rtr.Get("/messages", a.handleGetMessages)
	a.rtr.Get("/messages/:id", a.handleGetMessage)
	a.rtr.Delete("/messages/:id", a.handleDeleteMessage)

	return a
}

type ApiOption func(*api)

func WithMq(mq *mq.Mq) ApiOption {
	return func(a *api) {
		a.mq = mq
	}
}

func (a *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.rtr.ServeHTTP(w, r)
}

func (a *api) URLParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

func (a *api) render(w http.ResponseWriter, s int, d interface{}) {
	b, err := json.Marshal(d)
	if err != nil {
		http.Error(w, "failed to marshal data", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(s)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
