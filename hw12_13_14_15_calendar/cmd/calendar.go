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

// calendarCmd represents the calendar command.
var calendarCmd = &cobra.Command{
	Use:   "calendar",
	Short: "Calendar gRPC API",
	Long: `Calendar provides methods and abstraction calls
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
				log.Fatal(err)
				close(s)
			}
		}()
		go func() {
			err := server.RunProxy()
			if err != nil {
				log.Fatal("cannot run REST API proxy: %w", err)
				close(s)
			}
		}()
		defer server.Close()

		// // Defined server and waiting for shutdown.
		// HTTPServer := internalhttp.NewHTTPServer(cfg, zapLogger, database)
		// defer func() {
		// 	if err := HTTPServer.Shutdown(context.Background()); err != nil {
		// 		zapLogger.Error(err.Error())
		// 	}
		// 	zapLogger.Info("HTTP server shutdowned")
		// }()

		// // Main server goroutine.
		// go func() {
		// 	err := HTTPServer.Run()
		// 	if err != nil {
		// 		zapLogger.Error(err.Error())

		// 		// Don't wait for signal on server error.
		// 		close(s)
		// 	}
		// }()

		// Wait for interrupt signals.
		<-s
	},
}

func init() {
	rootCmd.AddCommand(calendarCmd)
}
