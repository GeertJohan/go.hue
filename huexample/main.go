package main

import (
	"fmt"
	"github.com/GeertJohan/go.hue"
	"github.com/davecgh/go-spew/spew"
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
	spew.Dump(brokerDetailsSlice)

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

	bridgeConfiguration, err := bridge.FetchConfiguration()
	if err != nil {
		fmt.Printf("have error: %s\n", err)
		return
	}

	fmt.Printf("Bridge name is '%s'.\n", bridgeConfiguration.Name)
	fmt.Printf("Bridge has %d users.\n", len(bridgeConfiguration.Whitelist))

	lights, err := bridge.Lights()
	if err != nil {
		fmt.Printf("have error: %s\n", err)
		return
	}
	fmt.Printf("Have %d lights.\n", len(lights))

	spew.Dump(lights[0].Attributes())

	// err = lights[0].SetName("Lange lamp")
	// if err != nil {
	// 	fmt.Printf("Have error: %s\n", err)
	// }

}
