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

	drone := &tello.Tello{}

	em.Println("✈️", " Preparing the flight...")
	err := drone.ControlConnectDefault()
	if err != nil {
		em.Printf("❌", "Failed to establish connection: %s\n", err)
		os.Exit(1)
	}

	em.Println("🛫", "Starting a quick journey...")
	drone.TakeOff()
	time.Sleep(5 * time.Second)

	em.Println("⏩", "Flying left...")
	move(drone, "left")
	time.Sleep(2 * time.Second)

	em.Println("⏪", "Flying Right...")
	move(drone, "right")
	time.Sleep(2 * time.Second)

	em.Println("🛬", "Landing...")
	drone.Land()
	drone.ControlDisconnect()

	em.Println("🏡", "Welcome home!")
}

func move(drone *tello.Tello, move string) {
	switch move {
	case "left":
		drone.Left(80)
	case "right":
		drone.Right(80)
	}

	time.Sleep(500 * time.Millisecond)
	drone.Hover()
}
