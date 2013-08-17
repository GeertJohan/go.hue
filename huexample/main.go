package main

import (
	"fmt"
	"github.com/GeertJohan/go.hue"
)

func main() {
	fmt.Println("Welcome to huexample")

	brokerDetailsSlice, err := hue.DiscoverBridges()
	if err != nil {
		fmt.Println("Error while discovering bridges:", err)
		return
	}

	if len(brokerDetailsSlice) == 0 {
		fmt.Println("No bridge details found. Stopping.")
		return
	}

	fmt.Printf("Found %d bridges:\n", len(brokerDetailsSlice))

	for idx, brokerDetails := range brokerDetailsSlice {
		fmt.Printf("%d: %#v\n", idx, brokerDetails)
	}

	fmt.Printf("Continueing with first bridge found, id: %s\n", brokerDetailsSlice[0].ID)

	fmt.Println("Going to create user with empty username, bridge will generate a username.")
	bridge := hue.NewBridge(brokerDetailsSlice[0].InternalIPAddress)
	newUsername, err := bridge.CreateNewUser("huexample", "")
	if err != nil {
		fmt.Printf("have error: %s\n", err)
		return
	}
	fmt.Printf("Successfully created new user. Got username: %s\n", newUsername)

	// update the Username field on Bridge instance with the user we just created
	bridge.Username = newUsername
}
