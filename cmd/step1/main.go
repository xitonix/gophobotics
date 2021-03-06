package main

import (
	"os"
	"time"

	"github.com/SMerrony/tello"
	"github.com/spf13/pflag"
	"github.com/xitonix/gophobotics/prnt"
)

func main() {
	de := pflag.BoolP("disable-emoticons", "d", false, "Disables emoticon printing")
	pflag.Parse()
	em := prnt.NewEmotifier(!*de)

	drone := tello.Tello{}

	em.Println("✈️", " Preparing the flight...")
	err := drone.ControlConnectDefault()
	if err != nil {
		em.Printf("❌", "Failed to establish connection: %s\n", err)
		os.Exit(1)
	}
	em.Println("🛫", "Starting a 10 seconds journey...")
	drone.TakeOff()
	time.Sleep(10 * time.Second)
	em.Println("🛬", "Landing...")
	drone.Land()
	drone.ControlDisconnect()
	em.Println("🏡", "Welcome home!")
}
