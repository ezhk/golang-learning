package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	localTime := time.Now()
	ntpTime, err := ntp.Time("pool.ntp.org")
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println("current time:", localTime.Round(0))
	fmt.Println("exact time:", ntpTime.Round(0))
}
