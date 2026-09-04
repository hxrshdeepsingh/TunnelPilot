package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	stdruntime "runtime"
	"sync"
	"syscall"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type TunnelService struct {
	event *application.EventManager
	mu    sync.Mutex
	procs map[int]*exec.Cmd
}

func (s *TunnelService) RunCommands(cmd string) string {
	cloudflared, err := getCloudflaredPath()
	if err != nil {
		return err.Error()
	}

	out, err := exec.Command(cloudflared, cmd).CombinedOutput()
	if err != nil {
		return err.Error()
	}

	return string(out)
}

func (s *TunnelService) RunTunnelLogin() string {
	cloudflared, err := getCloudflaredPath()
	if err != nil {
		return err.Error()
	}

	out, err := exec.Command(
		cloudflared,
		"tunnel",
		"login",
	).CombinedOutput()

	if err != nil {
		return err.Error()
	}

	return string(out)
}

func getCloudflaredPath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}

	appDir := filepath.Dir(exe)

	var path string

	if stdruntime.GOOS == "windows" {
		path = filepath.Join(appDir, "cloudflared.exe")
	} else {
		path = filepath.Join(appDir, "cloudflared")
	}

	if _, err := os.Stat(path); err != nil {
		return "", fmt.Errorf(
			"cloudflared not found at %s",
			path,
		)
	}

	return path, nil
}

func (s *TunnelService) RunTunnelLogout() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "Failed to get home directory: " + err.Error()
	}

	certPath := filepath.Join(
		home,
		".cloudflared",
		"cert.pem",
	)

	if err := os.Remove(certPath); err != nil {
		if os.IsNotExist(err) {
			return "Already logged out (no certificate found)."
		}

		return "Failed to logout: " + err.Error()
	}

	return "Successfully logged out. Certificate deleted."
}

func (s *TunnelService) CheckIsLoggedIn() bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	certPath := filepath.Join(
		home,
		".cloudflared",
		"cert.pem",
	)

	_, err = os.Stat(certPath)

	return err == nil
}

func (s *TunnelService) RunTunnelCreate(port string) int {
	cloudflared, err := getCloudflaredPath()
	if err != nil {
		s.event.Emit(
			"tunnel-log",
			err.Error(),
		)

		return 0
	}

	cmd := exec.Command(
		cloudflared,
		"tunnel",
		"--url",
		"http://localhost:"+port,
	)

	// Hide the console window on Windows.
	if stdruntime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow: true,
		}
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		s.event.Emit(
			"tunnel-log",
			"Failed to create stdout pipe: "+err.Error(),
		)

		return 0
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		s.event.Emit(
			"tunnel-log",
			"Failed to create stderr pipe: "+err.Error(),
		)

		return 0
	}

	if err := cmd.Start(); err != nil {
		s.event.Emit(
			"tunnel-log",
			"Start error: "+err.Error(),
		)

		return 0
	}

	pid := cmd.Process.Pid

	s.mu.Lock()
	s.procs[pid] = cmd
	s.mu.Unlock()

	emit := func(r io.Reader, tag string) {
		sc := bufio.NewScanner(r)

		sc.Buffer(
			make([]byte, 0, 64*1024),
			1024*1024,
		)

		for sc.Scan() {
			s.event.Emit(
				"tunnel-log",
				fmt.Sprintf(
					"[%d %s] %s",
					pid,
					tag,
					sc.Text(),
				),
			)
		}
	}

	go emit(stdout, "stdout")
	go emit(stderr, "stderr")

	go func() {
		_ = cmd.Wait()

		s.mu.Lock()
		delete(s.procs, pid)
		s.mu.Unlock()

		s.event.Emit(
			"tunnel-log",
			fmt.Sprintf("[%d] exited", pid),
		)
	}()

	return pid
}

func (s *TunnelService) StopAnyPID(pid int) bool {
	p, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	if err := p.Kill(); err != nil {
		return false
	}

	return true
}

func (s *TunnelService) KillAllTunnels() string {
	var cmd *exec.Cmd

	if stdruntime.GOOS == "windows" {
		cmd = exec.Command(
			"taskkill",
			"/IM",
			"cloudflared.exe",
			"/F",
		)
	} else {
		cmd = exec.Command(
			"killall",
			"-9",
			"cloudflared",
		)
	}

	out, err := cmd.CombinedOutput()

	s.mu.Lock()
	s.procs = make(map[int]*exec.Cmd)
	s.mu.Unlock()

	if err != nil {
		return fmt.Sprintf(
			"Finished kill attempt.\nOutput: %s",
			string(out),
		)
	}

	return fmt.Sprintf(
		"Successfully killed all cloudflared processes.\nOutput: %s",
		string(out),
	)
}
