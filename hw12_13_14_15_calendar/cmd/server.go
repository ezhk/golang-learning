// +build server

package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	logger "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/server/grpc"
	"github.com/spf13/cobra"
)

// serverCmd represents the calendar command.
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Calendar gRPC and REST API",
	Long: `Calendar server provides methods and abstraction calls
under the hood, that processing request as SQL commands.`,
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

		// Caught system signals.
		s := make(chan os.Signal)
		signal.Notify(s, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		// Run server.
		server := internalgrpc.NewServer(cfg, zapLogger, database)
		go func() {
			err := server.RunServer()
			if err != nil {
				zapLogger.Sugar().Errorf("cannot run gRPC server: %s", err)
				// Close program with closed channel.
				close(s)
			}
		}()
		go func() {
			err := server.RunProxy()
			if err != nil {
				zapLogger.Sugar().Errorf("cannot run REST API proxy: %s", err)
				close(s)
			}
		}()
		defer server.Close()

		// Wait for interrupt signals.
		<-s
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
