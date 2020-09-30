package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	logger "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/server/http"
	"github.com/spf13/cobra"
)

// calendarCmd represents the calendar command.
var calendarCmd = &cobra.Command{
	Use:   "calendar",
	Short: "Calendar gRPC API",
	Long: `Calendar provides methods and abstraction calls
under the hood, that processing request as SQL commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Init logger.
		log := logger.NewLogger(cfg)
		defer log.Close()

		// Init connect to database on runiing app.
		DSNString := cfg.GetDatabasePath()
		database := cfg.DatabaseBuilder()
		err := database.Connect(DSNString)
		if err != nil {
			log.Fatal("cannot conect to database: %w", err)
		}
		defer database.Close()

		// Caught system signals.
		s := make(chan os.Signal)
		signal.Notify(s, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		// Defined server and waiting for shutdown.
		HTTPServer := internalhttp.NewHTTPServer(cfg, log, database)
		defer func() {
			if err := HTTPServer.Shutdown(context.Background()); err != nil {
				log.Error(err)
			}
			log.Info("HTTP server shutdowned")
		}()

		// Main server goroutine.
		go func() {
			err := HTTPServer.Run()
			if err != nil {
				log.Error(err)

				// Don't wait for signal on server error.
				close(s)
			}
		}()

		// Wait for interrupt signals.
		<-s
	},
}

func init() {
	rootCmd.AddCommand(calendarCmd)
}
