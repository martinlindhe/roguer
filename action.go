package rogue

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
)

// Action ...
type Action interface {
	// returns true when finished performing action
	Perform(npc *Obj) bool
}

type actionList struct {
	All []actionSpec `json:"all"`
}

type actionSpec struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Result      string `json:"result"`
	Duration    int    `json:"duration"`
	Energy      int    `json:"energy"`
	TimeSpent   int
	Destination *Point
}

func parseActionsDefinition(defFileName string) ([]actionSpec, error) {

	var specs actionList

	data, err := ioutil.ReadFile(defFileName)
	if err != nil {
		return specs.All, err
	}

	err = yaml.Unmarshal(data, &specs)
	if err != nil {
		return specs.All, err
	}

	generalLog.Infof("Read", len(specs.All), "entries from", defFileName)
	return specs.All, nil
}
