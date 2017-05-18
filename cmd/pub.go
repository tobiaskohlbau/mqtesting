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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tobiaskohlbau/mqtesting/mqtt"
)

var (
	msg   string
	topic string
)

// pubCmd represents the pub command
var pubCmd = &cobra.Command{
	Use:   "pub",
	Short: "Publishes message on broker localhost:1883 use -m for message.",
	Run: func(cmd *cobra.Command, args []string) {
		if msg == "" || topic == "" {
			cmd.Usage()
			return
		}
		brokerAddress := viper.GetString("broker")
		if brokerAddress == "localhost:1883" {
			log.Println("pub: using default broker settings")
		}
		mqttService, err := mqtt.New(brokerAddress, "mqtesting-pub")
		if err != nil {
			log.Fatal(err)
		}
		if err := mqttService.Publish(topic, 0, false, msg); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(pubCmd)
	pubCmd.Flags().StringVarP(&msg, "message", "m", "", "message to send")
	pubCmd.Flags().StringVarP(&topic, "topic", "t", "", "topic to publish on")
}
