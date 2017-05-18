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
	"net/http"
	"strconv"
)

func (s *service) handleGetMessages(w http.ResponseWriter, r *http.Request) {
	tp := r.URL.Query().Get("topic")
	msgs, err := s.store.Messages().ListByTopic(tp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(msgs) == 0 {
		s.render(w, http.StatusOK, []string{})
		return
	}
	s.render(w, http.StatusOK, msgs)
}

func (s *service) handleGetMessage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(s.urlParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "failed to parse message id", http.StatusBadRequest)
		return
	}
	msg, err := s.store.Messages().Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.render(w, http.StatusOK, msg)
}

func (s *service) handleDeleteMessage(w http.ResponseWriter, r *http.Request) {
	usr := s.user(r)
	if usr == nil {
		w.Header().Set("WWW-Authenticate", "Bearer")
		http.Error(w, "authentication required", http.StatusUnauthorized)
		return
	}

	id, err := strconv.ParseInt(s.urlParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "failed to parse message id", http.StatusBadRequest)
		return
	}
	if err := s.store.Messages().Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
