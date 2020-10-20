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
	"github.com/ezhk/golang-learning/banners-rotation/internal/server"
	"github.com/ezhk/golang-learning/banners-rotation/internal/storage"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command.
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "gRPC and HTTP server",
	Long: `Server providers gRPC and REST operations
under exist objects like a banners, slots and groups.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Init logger.
		zapLogger := logger.NewLogger(cfg)
		defer zapLogger.Close()

		// Init connect to database on runiing app.
		storage, err := storage.NewStorage(cfg)
		if err != nil {
			log.Fatal("cannot create new storage: %w", err)
		}

		server := server.NewServer(cfg, zapLogger, storage)

		// Caught system signals.
		s := make(chan os.Signal)
		signal.Notify(s, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		// Run gRPC server.
		go func() {
			err := server.Run()
			if err != nil {
				close(s)
				log.Fatal(err)
			}
		}()

		// Run REST API HTTP server.
		go func() {
			err := server.RunProxy()
			if err != nil {
				log.Fatal("cannot run REST API proxy: %w", err)
				close(s)
			}
		}()

		// Wait for interrupt signals.
		<-s
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
