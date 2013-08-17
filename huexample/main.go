package main

import (
	"fmt"
	"github.com/GeertJohan/go.hue"
)

func main() {
	fmt.Println("Welcome to huexamlpe")

	bbds, err := hue.DiscoverBridges()
	if err != nil {
		fmt.Println("Error while discovering bridges:", err)
		return
	}

	if len(bbds) == 0 {
		fmt.Println("No bridge details found. Stopping.")
		return
	}

	fmt.Printf("Found %d bridges:\n", len(bbds))

	for idx, bbd := range bbds {
		fmt.Printf("%d: %#v\n", idx, bbd)
	}

	fmt.Printf("Continueing with first bridge found, id: %s\n", bbds[0].ID)

	fmt.Println("Going to create user 'newdeveloper':")
	bridge1 := hue.NewBridge(bbds[0].InternalIPAddress)
	newUsername, err := bridge1.CreateNewUser("huexample", "newdeveloper")
	if err != nil {
		fmt.Printf("have error: %s\n", err)
		return
	}
	fmt.Printf("Successfully created new user. Got username: %s\n", newUsername)

	fmt.Println("Going to create user with empty username:")
	bridge2 := hue.NewBridge(bbds[0].InternalIPAddress)
	newUsername, err = bridge2.CreateNewUser("huexample", "")
	if err != nil {
		fmt.Printf("have error: %s\n", err)
		return
	}
	fmt.Printf("Successfully created new user. Got username: %s\n", newUsername)

	bridge2.Username = newUsername
}
