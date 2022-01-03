package cfg

import (
	"encoding/json"
	"io/ioutil"
)

var Config Configuration

type Configuration struct {
	Port         int         `json:"Port"`
	MIDIPortName string      `json:"MIDIPortName"`
	TrackCoords  [][]float64 `json:"TrackCoords"`
}

func LoadConfig() error {
	cfgFile, err := ioutil.ReadFile("config.json")
	if err != nil {
		return err
	}

	if err := json.Unmarshal(cfgFile, &Config); err != nil {
		return err
	}

	return nil
}
