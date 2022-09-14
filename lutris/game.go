package lutris

import (
	"errors"
	"fmt"
	"os/exec"
	"syscall"
)

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
