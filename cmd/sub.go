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
	"log"
	"net/http"

	"github.com/pressly/chi"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cobra"
	"github.com/tobiaskohlbau/mqtesting/api"
	"github.com/tobiaskohlbau/mqtesting/mq"
	_ "github.com/tobiaskohlbau/mqtesting/statik"
	"github.com/tobiaskohlbau/mqtesting/store/postgres"
)

var (
	dbAddr     string
	dbUser     string
	dbPassword string
	dbName     string
)

// subCmd represents the sub command
var subCmd = &cobra.Command{
	Use:   "sub",
	Short: "Subscribes to specified address with all topics and serves on http://localhost the latest received messages",
	Run: func(cmd *cobra.Command, args []string) {
		str, err := postgres.Connect(dbAddr, dbUser, dbPassword, dbName)
		if err != nil {
			log.Fatal(err)
		}
		mq := mq.New(mqttAddress, str)
		a := api.New(str, api.WithMq(mq))

		r := chi.NewRouter()

		r.Mount("/api", a)

		statikFS, _ := fs.New()
		r.FileServer("/", statikFS)

		http.ListenAndServe(":http", r)
	},
}

func init() {
	RootCmd.AddCommand(subCmd)
	subCmd.Flags().StringVar(&dbAddr, "dbAddress", "localhost:5432", "db host:port")
	subCmd.Flags().StringVar(&dbUser, "dbUser", "postgres", "db user")
	subCmd.Flags().StringVar(&dbPassword, "dbPassword", "secret", "db password")
	subCmd.Flags().StringVar(&dbName, "dbName", "postgres", "db name")
}
