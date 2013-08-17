package hue

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type apiResponse struct {
	Success map[string]string
	Error   *apiResponseError `json:"error"`
}

type apiResponseError struct {
	Type        uint   `json:"type"`
	Address     string `json:"address"`
	Description string `json:"description"`
}

// Bridge represents a Hue Bridge
type Bridge struct {
	id       string
	IP       string
	Username string
}

// NewBridge creates a new Bridge instance with given IP address
func NewBridge(IP string) *Bridge {
	b := &Bridge{
		IP: IP,
	}
	return b
}

// ID returns the ID of the Bridge as string
func (b *Bridge) ID() string {
	if len(b.id) == 0 {
		//++ retrieve ID
	}
	return b.id
}

func (b *Bridge) URL() string {
	return "http://" + b.IP + "/api"
}

// CreateNewUser creates a new user at the bridge.
// The end-user must press the link button in advance to prove physical access.
// When the second argument (newUsername) is left emtpy, the bridge will provide a username.
// CreateNewUser does not update the Bridge instance with the username. This must be done manually.
func (b *Bridge) CreateNewUser(deviceType string, newUsername string) (string, error) {
	requestData := map[string]string{"devicetype": deviceType}
	if len(newUsername) > 0 {
		requestData["username"] = newUsername
	}

	// create empty buffer
	buf := bytes.NewBuffer(nil)

	// encode requestData to buffer
	err := json.NewEncoder(buf).Encode(requestData)
	if err != nil {
		return "", err
	}

	// do post to api
	response, err := http.Post(b.URL(), "text/json", buf)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// create and decode apiResponse
	apiResponseSlice := make([]*apiResponse, 0, 1)
	err = json.NewDecoder(response.Body).Decode(&apiResponseSlice)
	if err != nil {
		return "", err
	}
	if len(apiResponseSlice) == 0 {
		return "", errors.New("received empty api response array")
	}
	if len(apiResponseSlice) > 1 {
		return "", errors.New("received api response array with >1 items")
	}

	apiResponse := apiResponseSlice[0]

	// check for error from bridge
	if apiResponse.Error != nil {
		return "", errors.New(apiResponse.Error.Description)
	}

	return apiResponse.Success["username"], nil
}
