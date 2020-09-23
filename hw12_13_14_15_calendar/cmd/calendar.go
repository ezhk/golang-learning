package cmd

import (
	logger "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/internal/server/http"
	"github.com/spf13/cobra"
)

// calendarCmd represents the calendar command.
var calendarCmd = &cobra.Command{
	Use:   "calendar",
	Short: "Calendar CLI",
	Long: `Calendar provides methods and abstraction calls
under the hood, that processing request as SQL commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		log := logger.NewLogger()
		log.Info("Calendar has called")

		HTTPServer := internalhttp.NewHTTPServer(log)
		err := HTTPServer.Run()
		if err != nil {
			log.Error("Received HTTP run error", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(calendarCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// calendarCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// calendarCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
