package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	stdruntime "runtime"
	"sync"
	"runtime"
	// "syscall"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type TunnelService struct {
	event *application.EventManager
	mu    sync.Mutex
	procs map[int]*exec.Cmd
}

// func newHiddenCommand(name string, args ...string) *exec.Cmd {
// 	cmd := exec.Command(name, args...)

// 	if stdruntime.GOOS == "windows" {
// 		cmd.SysProcAttr = &syscall.SysProcAttr{
// 			HideWindow:    true,
// 			CreationFlags: 0x08000000,
// 		}
// 	}

// 	return cmd
// }

func newHiddenCommand(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	configureHiddenCommand(cmd)
	return cmd
}

func (s *TunnelService) RunCommands(cmd string) string {
	cloudflared, err := ensureCloudflared()
	if err != nil {
		return err.Error()
	}

	out, err := newHiddenCommand(cloudflared, cmd).CombinedOutput()
	if err != nil {
		return err.Error()
	}

	return string(out)
}

func (s *TunnelService) RunTunnelLogin() string {
	cloudflared, err := ensureCloudflared()
	if err != nil {
		return err.Error()
	}

	out, err := newHiddenCommand(
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
	base, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	runtimeDir := filepath.Join(
		base,
		"TunnelPilot",
		"runtime",
	)

	if err := os.MkdirAll(runtimeDir, 0755); err != nil {
		return "", fmt.Errorf(
			"failed to create runtime directory: %w",
			err,
		)
	}

	return filepath.Join(runtimeDir, "cloudflared.exe"), nil
}

// func ensureCloudflared() (string, error) {
// 	path, err := getCloudflaredPath()
// 	if err != nil {
// 		return "", err
// 	}

// 	// Already downloaded.
// 	if _, err := os.Stat(path); err == nil {
// 		return path, nil
// 	}

// 	const url = "https://github.com/cloudflare/cloudflared/releases/download/2026.9.3/cloudflared-windows-amd64.exe"

// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return "", fmt.Errorf(
// 			"failed to download cloudflared: %w",
// 			err,
// 		)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return "", fmt.Errorf(
// 			"cloudflared download returned HTTP %d",
// 			resp.StatusCode,
// 		)
// 	}

// 	file, err := os.Create(path)
// 	if err != nil {
// 		return "", fmt.Errorf(
// 			"failed to create cloudflared: %w",
// 			err,
// 		)
// 	}
// 	defer file.Close()

// 	if _, err := io.Copy(file, resp.Body); err != nil {
// 		return "", fmt.Errorf(
// 			"failed to save cloudflared: %w",
// 			err,
// 		)
// 	}

// 	return path, nil
// }


func ensureCloudflared() (string, error) {
	path, err := getCloudflaredPath()
	if err != nil {
		return "", err
	}

	// Select the correct release asset for the current OS and architecture.
	var url string

	switch runtime.GOOS {
	case "windows":
		if runtime.GOARCH != "amd64" {
			return "", fmt.Errorf("unsupported Windows architecture: %s", runtime.GOARCH)
		}
		url = "https://github.com/cloudflare/cloudflared/releases/download/2026.9.3/cloudflared-windows-amd64.exe"

	case "linux":
		switch runtime.GOARCH {
		case "amd64":
			url = "https://github.com/cloudflare/cloudflared/releases/download/2026.9.3/cloudflared-linux-amd64"
		case "arm64":
			url = "https://github.com/cloudflare/cloudflared/releases/download/2026.9.3/cloudflared-linux-arm64"
		default:
			return "", fmt.Errorf("unsupported Linux architecture: %s", runtime.GOARCH)
		}

	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	// Reuse the cached binary only if it exists and is usable.
	if info, err := os.Stat(path); err == nil && !info.IsDir() {
		if runtime.GOOS != "linux" || info.Mode()&0111 != 0 {
			return path, nil
		}
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download cloudflared: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("cloudflared download returned HTTP %d", resp.StatusCode)
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create cloudflared: %w", err)
	}

	_, copyErr := io.Copy(file, resp.Body)
	closeErr := file.Close()

	if copyErr != nil {
		_ = os.Remove(path)
		return "", fmt.Errorf("failed to save cloudflared: %w", copyErr)
	}
	if closeErr != nil {
		_ = os.Remove(path)
		return "", fmt.Errorf("failed to close cloudflared: %w", closeErr)
	}

	// Linux needs executable permissions to launch the downloaded binary.
	if runtime.GOOS == "linux" {
		if err := os.Chmod(path, 0755); err != nil {
			_ = os.Remove(path)
			return "", fmt.Errorf("failed to set cloudflared executable permissions: %w", err)
		}
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
	cloudflared, err := ensureCloudflared()
	if err != nil {
		s.event.Emit(
			"tunnel-log",
			err.Error(),
		)

		return 0
	}

	cmd := newHiddenCommand(
		cloudflared,
		"tunnel",
		"--url",
		"http://localhost:"+port,
	)

	// Hide the console window on Windows.
	// if stdruntime.GOOS == "windows" {
	// 	cmd.SysProcAttr = &syscall.SysProcAttr{
	// 		HideWindow: true,
	// 	}
	// }

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
		cmd = newHiddenCommand(
			"taskkill",
			"/IM",
			"cloudflared.exe",
			"/F",
		)
	} else {
		cmd = newHiddenCommand(
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
