package state

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

type Package struct {
	Prefix  string   `json:"prefix"`
	Release string   `json:"release"`
	Local   bool     `json:"local"`
	Files   []string `json:"files"`
}

type State = map[string]Package

func statePath() (string, error) {
	dir := filepath.Join(xdg.StateHome, "tuck")
	info, _ := os.Stat(dir)
	if info == nil || !info.IsDir() {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return "", err
		}
	}
	return filepath.Join(dir, "installed.json"), nil
}

func LoadState() (State, error) {
	state := make(State)
	path, err := statePath()
	if err != nil {
		return state, err
	}
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

func DumpState(state State) error {
	path, err := statePath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func Install(name string, pkg Package) error {
	state, err := LoadState()
	if err != nil {
		return err
	}
	state[name] = pkg
	return DumpState(state)
}

func Remove(name string) error {
	state, err := LoadState()
	if err != nil {
		return err
	}
	_, found := state[name]
	if found {
		delete(state, name)
	}
	return DumpState(state)
}
