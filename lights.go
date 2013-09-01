package hue

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

// Light points to a specific light on a specific hue bridge
type Light struct {
	bridge *Bridge // bridge on which the light is connected
	id     string  // id of the light
}

func (l Light) Attributes() (*LightAttributes, error) {
	resp, err := http.Get(l.bridge.URL() + "/lights/" + l.id)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	attributes := &LightAttributes{}
	err = json.NewDecoder(resp.Body).Decode(attributes)
	if err != nil {
		return nil, err
	}

	return attributes, nil
}

// SetName sets the name of the light. The given name must have a length between 0 and 32 characters.
func (l Light) SetName(newName string) error {
	//++ TODO: check for ascii characters only??
	if len(newName) > 32 {
		return errors.New("given name exceeds length limit")
	}

	setNameObject := map[string]string{"name": newName}
	bodyBytes, err := json.Marshal(setNameObject)
	if err != nil {
		return err
	}
	bodyBuf := bytes.NewBuffer(bodyBytes)

	client := &http.Client{}
	request, err := http.NewRequest("PUT", l.bridge.URL()+"/lights/"+l.id+"/name", bodyBuf)
	if err != nil {
		return err
	}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// LightAttributes holds attributes of light, it includes the State and Name.
type LightAttributes struct {
	State     LightState `json:"State"`     // Details the state of the light, see the state table below for more details.
	Type      string     `json:"Type"`      // A fixed name describing the type of light e.g. “Extended color light”.
	Name      string     `json:"name"`      // (lenght 0-32) A unique, editable name given to the light.
	ModelID   string     `json:"modelid"`   // (length 6) The hardware model of the light.
	Swversion string     `json:"swversion"` // (length 8) An identifier for the software version running on the light.
	// Pointsymbol string     `json:"Pointsymbol"` // (object) This parameter is reserved for future functionality.
}

type LightState struct {
	On         bool   `json:"On"`  // On/Off state of the light. On=true, Off=false
	Brightness uint8  `json:"Bri"` // Brightness of the light. This is a scale from the minimum brightness the light is capable of, 0, to the maximum capable brightness, 255. Note a brightness of 0 is not off.
	Hue        uint16 `json:"Hue"` // Hue of the light. This is a wrapping value between 0 and 65535. Both 0 and 65535 are red, 25500 is green and 46920 is blue.
	Saturation uint8  `json:"sat"` // Saturation of the light. 255 is the most saturated (colored) and 0 is the least saturated (white).

	// xy ?? // 2..2 of float 4	The x and y coordinates of a color in CIE color space.
	// The first entry is the x coordinate and the second entry is the y coordinate. Both x and y are between 0 and 1.

	CT uint16 `json:"ct"` // The Mired Color temperature of the light. 2012 connected lights are capable of 153 (6500K) to 500 (2000K).

	Alert string `json:"alert"` // The alert effect, which is a temporary change to the bulb’s state. This can take one of the following values:
	// “none” – The light is not performing an alert effect.
	// “select” – The light is performing one breathe cycle.
	// “lselect” – The light is performing breathe cycles for 30 seconds or until an "alert": "none" command is received.
	// Note that in version 1.0 this contains the last alert sent to the light and not its current state. This will be changed to contain the current state in an upcoming patch.

	Effect string `json:"effect"` // The dynamic effect of the light, can either be “none” or “colorloop”.

	// If set to colorloop, the light will cycle through all hues using the current brightness and saturation settings.
	ColorMode string `json:"colormode"` // (length 2) Indicates the color mode in which the light is working, this is the last command type it received. Values are “hs” for Hue and Saturation, “xy” for XY and “ct” for Color Temperature. This parameter is only present when the light supports at least one of the values.
	Reachable bool   `json:"reachable"` // Indicates if a light can be reached by the bridge. Currently always returns true, functionality will be added in a future patch.
}

// Lights returns all lights known by the bridge.
func (b *Bridge) Lights() ([]Light, error) {
	resp, err := http.Get(b.URL() + "/lights")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	lightsMap := map[string]interface{}{} // we use interface{} to discard the value on each key
	err = json.NewDecoder(resp.Body).Decode(&lightsMap)
	if err != nil {
		return nil, err
	}
	lights := make([]Light, 0, len(lightsMap))
	for lightID, _ := range lightsMap {
		lights = append(lights, Light{b, lightID})
	}
	return lights, nil
}

// Search lets the bridge start a new search for lights.
// The bridge will search for 1 minute and will add a maximum of 15 new lights.
// To add further lights, the command needs to be sent again after the search has completed.
// If a search is already active, it will be aborted and a new search will start.
func (b *Bridge) Search() error {
	resp, err := http.Post(b.URL()+"/lights", "", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
