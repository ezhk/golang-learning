// +build scheduler

package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/broker"
	logger "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/logger"
	"github.com/spf13/cobra"
)

// schedulerCmd represents the scheduler command.
var schedulerCmd = &cobra.Command{
	Use:   "scheduler",
	Short: "Schedule messages for undelivered events",
	Long: `Process finds database events,
that occur two or less weeks, and generates
messages into broker (rabbitMQ).`,
	Run: func(cmd *cobra.Command, args []string) {
		// Init logger.
		zapLogger := logger.NewLogger(cfg)
		defer zapLogger.Close()

		// Init connect to database on runiing app.
		DSNString := cfg.GetDatabasePath()
		database := cfg.DatabaseBuilder()

		err := database.Connect(DSNString)
		if err != nil {
			log.Fatal("cannot conect to database: %w", err)
		}
		defer database.Close()

		// Init broker producer.
		producer := broker.NewProducer(cfg)
		err = producer.Init()
		if err != nil {
			log.Fatal("cannot initialize producer: %w", err)
		}
		defer producer.Close()

		// Caught system signals.
		s := make(chan os.Signal)
		signal.Notify(s, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			for {
				events, err := database.GetNotifyReadyEvents()
				if err != nil {
					zapLogger.Sugar().Errorf("cannot get notiry ready events: %s", err)
					// Close program with closed channel.
					close(s)
				}

				for _, event := range events {
					err = producer.Publish(event)
					if err != nil {
						zapLogger.Sugar().Errorf("cannot publish event: %s", err)
						close(s)
					}

					event.Notified = true
					err = database.UpdateEvent(event)
					if err != nil {
						zapLogger.Sugar().Errorf("cannot mark event as processed: %s", err)
						close(s)
					}

					zapLogger.Sugar().Infof("scheduled event: %+v", event)
				}

				time.Sleep(time.Duration(cfg.Scheduler.CheckInterval) * time.Second)
			}
		}()

		// Wait for interrupt signals.
		<-s
	},
}

func init() {
	rootCmd.AddCommand(schedulerCmd)
}
