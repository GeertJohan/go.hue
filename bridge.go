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

// TODO: use time.Time for Utc and Whitelist.LastUseDate, Whitelist.CreateDate
type BridgeConfiguration struct {
	Proxyport uint16 `json:"proxyport"` // Port of the proxy being used by the bridge. If set to 0 then a proxy is not being used.
	Utc       string `json:"utc"`       // Current time stored on the bridge.
	Name      string `json:"name"`      // length 4..16. Name of the bridge. This is also its uPnP name, so will reflect the actual uPnP name after any conflicts have been resolved.
	SwUpdate  struct {
		UpdateState int    `json:"updatestate"`
		Url         string `json:"url"`
		Text        string `json:"text"`
		Notify      bool   `json:"notify"`
	} `json:"swupdate"` // Contains information related to software updates.
	Whitelist map[string]struct {
		LastUseDate string `json:"last use date"`
		CreateDate  string `json:"create date"`
		Name        string `json:"name"`
	} `json:"whitelist"` // An array of whitelisted user IDs.
	Swversion      string `json:"swversion"`      // Software version of the bridge.
	ProxyAddress   string `json:"proxyaddress"`   // length 0..40. IP Address of the proxy server being used. A value of “none” indicates no proxy.
	Mac            string `json:"mac"`            // MAC address of the bridge.
	LinkButton     bool   `json:"linkbutton"`     // Indicates whether the link button has been pressed within the last 30 seconds.
	IPAddress      string `json:"ipaddress"`      // IP address of the bridge.
	Netmask        string `json:"netmask"`        // Network mask of the bridge.
	Gateway        string `json:"gateway"`        // Gateway IP address of the bridge.
	DHCP           bool   `json:"dhcp"`           // Whether the IP address of the bridge is obtained with DHCP.
	PortalServices bool   `json:"portalservices"` // This indicates whether the bridge is registered to synchronize data with a portal account.
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
	return "http://" + b.IP + "/api/" + b.Username
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
	response, err := http.Post("http://"+b.IP+"/api", "text/json", buf)
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

// FetchConfiguration fetches the configuration data and returns it as *BridgeConfiguration
func (b *Bridge) FetchConfiguration() (*BridgeConfiguration, error) {
	response, err := http.Get(b.URL() + "/config")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// create and decode apiResponse
	bridgeConfiguration := &BridgeConfiguration{}
	err = json.NewDecoder(response.Body).Decode(bridgeConfiguration)
	if err != nil {
		return nil, err
	}

	return bridgeConfiguration, nil
}
