package rogue

import (
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
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
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	Result    []string `json:"result"`
	Duration  int      `json:"duration"`
	Energy    int      `json:"energy"`
	TimeSpent int
}

func parseActionsDefinition(defFileName string) []actionSpec {

	data, err := ioutil.ReadFile(defFileName)
	if err != nil {
		panic(err)
	}

	var actions actionList
	err = yaml.Unmarshal(data, &actions)
	if err != nil {
		panic(err)
	}

	log.Infof("Read %d entries from %s", len(actions.All), defFileName)
	return actions.All
}
