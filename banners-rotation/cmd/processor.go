/*
Copyright © 2020 Andrey Kiselev <kiselevandrew@yandex.ru>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ezhk/golang-learning/banners-rotation/internal/logger"
	"github.com/ezhk/golang-learning/banners-rotation/internal/queue"
	"github.com/ezhk/golang-learning/banners-rotation/internal/storage"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command.
var processorCmd = &cobra.Command{
	Use:   "processor",
	Short: "event messages receiver",
	Long: `Server read events from redis and process them:
* update placement shows/clicks;
* calculate new score values.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Init logger.
		zapLogger := logger.NewLogger(cfg)
		defer zapLogger.Close()

		// Init connect to database on runiing app.
		storage, err := storage.NewStorage(cfg)
		if err != nil {
			log.Fatal("cannot create new storage: %w", err)
		}

		// Init connect to redis.
		queue, err := queue.NewQueue(cfg)
		if err != nil {
			log.Fatal("cannot create queue connect: %w", err)
		}

		err = queue.RunConsumer(storage, zapLogger)
		if err != nil {
			log.Fatal("cannot run consumer: %w", err)
		}

		// Caught system signals.
		s := make(chan os.Signal, 1)
		signal.Notify(s, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		// Wait for interrupt signals.
		<-s

		// Wait for consumer.
		<-queue.Conn.StopAllConsuming()
	},
}

func init() {
	rootCmd.AddCommand(processorCmd)
}