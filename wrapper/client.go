package wrapper

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

var LutrisExecutableNotFound = errors.New("lutris executable not found")
var LutrisCommandError = errors.New("lutris command failed")

type LutrisClient struct {
	lutrisPath string
}

func NewLutrisClient() (LutrisClient, error) {
	output, err := exec.Command("which", "lutris").Output()

	if err != nil {
		return LutrisClient{lutrisPath: ""}, LutrisExecutableNotFound
	}

	lutris_path := strings.Trim(string(output), " \n")

	return LutrisClient{lutrisPath: lutris_path}, nil
}

func (w *LutrisClient) FetchGames() ([]Game, error) {
	output, err := exec.Command(w.lutrisPath, "--list-games", "--json", "--installed").Output()

	if err != nil {
		return nil, errors.New(fmt.Sprintf("lutris command failed, error message: '%v'", output))
	}

	var games []Game

	err = json.Unmarshal(output, &games)

	if err != nil {
		return nil, err
	}

	for i := range games {
		games[i].lutrisPath = w.lutrisPath
	}

	return games, nil
}
