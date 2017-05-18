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
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type provider struct {
	cfg  *oauth2.Config
	user func(*http.Client) (string, string, error)
}

type providerConfig struct {
	endpoint oauth2.Endpoint
	scopes   []string
	user     func(*http.Client) (string, string, error)
}

var providerConfigs = map[string]providerConfig{
	"google": {
		endpoint: google.Endpoint,
		scopes:   []string{"profile"},
		user: func(c *http.Client) (string, string, error) {
			url := "https://www.googleapis.com/oauth2/v2/userinfo"
			u := struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}{}
			err := getJSON(c, url, &u)
			if err != nil {
				return "", "", err
			}
			return u.ID, u.Name, nil
		},
	},
	"microsoft": {
		endpoint: oauth2.Endpoint{
			AuthURL:  "https://login.microsoftonline.com/common/oauth2/v2.0/authorize",
			TokenURL: "https://login.microsoftonline.com/common/oauth2/v2.0/token",
		},
		scopes: []string{"https://graph.microsoft.com/user.read"},
		user: func(c *http.Client) (string, string, error) {
			url := "https://graph.microsoft.com/v1.0/users/me"
			u := struct {
				ID   string `json:"id"`
				Name string `json:"displayName"`
			}{}
			err := getJSON(c, url, &u)
			if err != nil {
				return "", "", err
			}
			return u.ID, u.Name, nil
		},
	},
}

func getJSON(c *http.Client, url string, v interface{}) error {
	r, err := c.Get(url)
	if err != nil {
		return fmt.Errorf("oauth request failed: %v", err)
	}
	defer r.Body.Close()

	if r.StatusCode < 200 || r.StatusCode > 299 {
		return fmt.Errorf("oauth wrong status code: %v", r.StatusCode)
	}

	if err := json.NewDecoder(io.LimitReader(r.Body, 1<<20)).Decode(v); err != nil {
		return fmt.Errorf("failed to decode JSON response: %v", err)
	}

	return nil
}
