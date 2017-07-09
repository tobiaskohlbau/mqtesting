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

package cmd

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"kohlbau.de/x/mqtesting/api"
	"kohlbau.de/x/mqtesting/jwt"
	"kohlbau.de/x/mqtesting/oauth"
	"kohlbau.de/x/mqtesting/store"
	"kohlbau.de/x/mqtesting/store/memory"
	"kohlbau.de/x/mqtesting/store/postgres"

	"net"

	"github.com/pressly/chi"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"kohlbau.de/x/mqtesting/mqtt"
	_ "kohlbau.de/x/mqtesting/statik"
)

const (
	defaultJWTSecret = "secret"
)

// subCmd represents the sub command
var subCmd = &cobra.Command{
	Use:   "sub",
	Short: "Subscribes to specified address with all topics and serves on http://localhost the latest received messages",
	Run: func(cmd *cobra.Command, args []string) {
		host, port, err := net.SplitHostPort(viper.GetString("address"))
		if err != nil {
			log.Fatal(err)
		}

		baseURL := viper.GetString("base_url")

		r := chi.NewRouter()

		var storeService store.Service
		dbConfig := viper.GetStringMapString("db")
		storeService, err = postgres.Connect(dbConfig["address"], dbConfig["user"], dbConfig["password"], dbConfig["name"])
		if err != nil {
			log.Println("sub: database config failed using in memory database")
			storeService = memory.New()
		}

		jwtSecret := viper.GetString("jwt")
		if jwtSecret == defaultJWTSecret {
			log.Println("sub: jwt secret not found using default value")
		}
		jwtService := jwt.New(jwt.WithSecret([]byte(jwtSecret)))

		oauthOptions := []oauth.Option{}
		for p, v := range viper.GetStringMapString("oauth") {
			s := strings.Split(v, ":")
			oauthOptions = append(oauthOptions, oauth.WithProvider(p, s[0], s[1]))
		}
		if url := viper.GetString("frontend_url"); url != "static" {
			oauthOptions = append(oauthOptions, oauth.WithFrontend(url))
		}
		oauthService := oauth.New(storeService.Users(), jwtService, fmt.Sprintf("%s/oauth", baseURL), oauthOptions...)
		r.Mount("/oauth", oauthService)

		brokerAddress := viper.GetString("broker")
		if brokerAddress == "localhost:1883" {
			log.Println("sub: using default broker settings")
		}
		mqttService, err := mqtt.New(brokerAddress, "mqtesting-sub")
		if err != nil {
			log.Fatal(err)
		}

		err = mqttService.Subscribe("#", func(c mqtt.Client, m mqtt.Message) {
			onMessage(storeService.Messages(), c, m)
		})
		if err != nil {
			log.Fatal(err)
		}

		apiOptions := []api.Option{
			api.WithStore(storeService),
			api.WithOAuth(oauthService),
			api.WithJWT(jwtService),
		}
		apiService := api.New(apiOptions...)
		r.Mount("/api", apiService)

		if viper.GetString("frontend_url") == "static" {
			statikFS, _ := fs.New()
			r.FileServer("/", statikFS)
		}

		if p := viper.GetInt("http_redirect"); p != 0 {
			go http.ListenAndServe(fmt.Sprintf("%s:%d", host, p), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, "https://"+r.Host+r.URL.String(), http.StatusMovedPermanently)
			}))
		}
		http.ListenAndServeTLS(fmt.Sprintf("%s:%s", host, port), viper.GetString("cert"), viper.GetString("key"), r)
	},
}

func onMessage(s store.MessageService, c mqtt.Client, msg mqtt.Message) {
	_, err := s.New(msg.Topic(), string(msg.Payload()[:]), msg.MessageID())
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	RootCmd.AddCommand(subCmd)

	subCmd.Flags().String("address", ":https", "address to listen on")
	subCmd.Flags().String("base_url", "https://localhost", "specifies base url")
	subCmd.Flags().Int("http_redirect", 80, "redirects http on specified port to https. 0 disables option")
	subCmd.Flags().String("frontend_url", "static", "specifies the frontend host:port by default static compiled in frontend")
	subCmd.Flags().String("cert", "certs/server.pem", "specifies server cert file")
	subCmd.Flags().String("key", "certs/server.key", "specifies server key file")
	subCmd.Flags().String("jwt", defaultJWTSecret, "specifies jwt secret")
	viper.BindPFlag("address", subCmd.Flags().Lookup("address"))
	viper.BindPFlag("base_url", subCmd.Flags().Lookup("base_url"))
	viper.BindPFlag("http_redirect", subCmd.Flags().Lookup("http_redirect"))
	viper.BindPFlag("host", subCmd.Flags().Lookup("host"))
	viper.BindPFlag("frontend_url", subCmd.Flags().Lookup("frontend_url"))
	viper.BindPFlag("cert", subCmd.Flags().Lookup("cert"))
	viper.BindPFlag("key", subCmd.Flags().Lookup("key"))
	viper.BindPFlag("jwt", subCmd.Flags().Lookup("jwt"))
}
