/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"github.com/spf13/cobra"

	database "../database"
	logger "../logger"
)

// calendarCmd represents the calendar command
var calendarCmd = &cobra.Command{
	Use:   "calendar",
	Short: "Calendar CLI",
	Long: `Calendar provides methods and abstraction calls
under the hood, that processing request as SQL commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("calendar called:", viper.Get("db.path"), viper.Get("logger"))

		log := logger.NewLogger()
		// log.Errorf("test %#v", viper.Get("db"))
		// fmt.Println(log.Close())

		db := database.NewDatatabase()
		log.Info(db.CreateUser("Vasya", "Pupkin"))
		log.Info(db.CreateUser("Vinni", "Pooh"))

		vinniPooh := db.GetUser("Vinni", "Pooh")
		log.Infof("%#v %v", db, vinniPooh)
		log.Info(db.UpdateUser(vinniPooh.Id, "Винни", "Пух"))
		log.Infof("%#v", db)

		db.DeleteUser(vinniPooh.Id)
		log.Infof("%#v", db)
		log.Info(db.CreateUser("Vasya", "Pupkin"))

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