// +build sender

package cmd

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/broker"
	logger "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/logger"
	"github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/structs"
	"github.com/spf13/cobra"
)

// senderCmd represents the sender command.
var senderCmd = &cobra.Command{
	Use:   "sender",
	Short: "Send notifications for events",
	Long: `Process reads rabbitMQ queries
parses them into event structure and deliver.
Deliver is print to STDOUT.`,
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup

		// Init logger.
		zapLogger := logger.NewLogger(cfg)
		defer zapLogger.Close()

		// Init broker producer.
		consumer := broker.NewConsumer(cfg)
		err := consumer.Init()
		if err != nil {
			log.Fatal("cannot initialize consumer: %w", err)
		}
		defer consumer.Close()

		// Caught system signals.
		s := make(chan os.Signal)
		signal.Notify(s, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		// Run reader events goroutine.
		wg.Add(1)
		go func() {
			defer wg.Done()

			events, err := consumer.Consume()
			if err != nil {
				zapLogger.Sugar().Errorf("cannot consume events: %s", err)
				close(s)
			}

			for {
				select {
				case <-s:
					return
				case msg, ok := <-events:
					if !ok {
						zapLogger.Sugar().Errorf("received EOF")

						return
					}

					var recvEvent structs.Event
					err = json.Unmarshal(msg.Body, &recvEvent)
					if err != nil {
						zapLogger.Sugar().Errorf("cannot read event: %s", err)
						err = msg.Nack(false, false)
						if err != nil {
							zapLogger.Sugar().Errorf("cannot negative ack message: %s", err)
						}
					}

					zapLogger.Sugar().Infof("Received message: %+v", recvEvent)
					err = msg.Ack(false)
					if err != nil {
						zapLogger.Sugar().Errorf("cannot Ack message: %s", err)
					}
				}
			}
		}()
		wg.Wait()
	},
}

func init() {
	rootCmd.AddCommand(senderCmd)
}
