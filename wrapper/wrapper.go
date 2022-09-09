package wrapper

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

var LutrisExecutableNotFound = errors.New("lutris executable not found")
var LutrisCommandError = errors.New("lutris command failed")

type Wrapper struct {
	lutrisPath string
}

func NewWrapper() (Wrapper, error) {
	output, err := exec.Command("which", "lutris").Output()

	if err != nil {
		return Wrapper{lutrisPath: ""}, LutrisExecutableNotFound
	}

	lutris_path := strings.Trim(string(output), " \n")

	return Wrapper{lutrisPath: lutris_path}, nil
}

func (w *Wrapper) FetchGames() ([]Game, error) {
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

type Game struct {
	Id         uint32 `json:"id"`
	Slug       string `json:"slug"`
	Name       string `json:"name"`
	Runner     string `json:"runner"`
	Platform   string `json:"platform"`
	Year       uint32 `json:"year"`
	Directory  string `json:"directory"`
	Hidden     bool   `json:"hidden"`
	PlayTime   string `json:"playtime"`
	LastPlayed string `json:"lastplayed"`
	IsRunning  bool
	lutrisPath string
	command    *exec.Cmd
}

func (g *Game) Start() (*exec.Cmd, error) {
	command := exec.Command(g.lutrisPath, fmt.Sprintf("lutris:rungameid/%v", g.Id))

	// Needed so that we can kill all children/grandchildren processes
	// by making them have the same PGID
	command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	err := command.Start()

	if err != nil {
		return nil, errors.New(fmt.Sprintf("lutris command failed, error: '%v'", err))
	}

	g.command = command

	return command, nil
}

func (g *Game) Stop() {
	syscall.Kill(-g.command.Process.Pid, syscall.SIGKILL)
	g.IsRunning = false
}
