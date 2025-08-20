package state

import (
	"encoding/json"
	"os"
	"path/filepath"
	"tuck/internal/path"
)

type Package struct {
	Prefix  string   `json:"prefix"`
	Release string   `json:"release"`
	Local   bool     `json:"local"`
	Files   []string `json:"files"`
}

type State = map[string]Package

func statePath() string {
	return filepath.Join(path.StateDir, "installed.json")
}

func load() (State, error) {
	state := make(State)
	path := statePath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return state, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return state, err
	}
	err = json.Unmarshal(data, &state)
	return state, err
}

func store(state State) error {
	path := statePath()
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func Install(name string, pkg Package) error {
	state, err := load()
	if err != nil {
		return err
	}
	state[name] = pkg
	return store(state)
}

func GetAll() (*State, error) {
	state, err := load()
	if err != nil {
		return nil, err
	}
	return &state, nil
}

func Get(name string) (*Package, error) {
	state, err := load()
	if err != nil {
		return nil, err
	}
	pkg, found := state[name]
	if !found {
		return nil, nil
	}
	return &pkg, nil
}

func Remove(name string) error {
	state, err := load()
	if err != nil {
		return err
	}
	_, found := state[name]
	if found {
		delete(state, name)
	}
	return store(state)
}
