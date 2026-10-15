//go:build linux

package main

import "os/exec"

func configureHiddenCommand(cmd *exec.Cmd) {
	// Linux does not need Windows console-window flags.
}
