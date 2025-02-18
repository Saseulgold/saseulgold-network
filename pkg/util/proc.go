package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

func ProcessIsRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	// Send signal to process to check if it's running
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

func ProcessStart(rootDir, processName string, pid int) error {
	pidFile := GetPIDFile(rootDir, processName)

	procDir := filepath.Join(rootDir, "proc")
	if err := os.MkdirAll(procDir, 0755); err != nil {
		return fmt.Errorf("failed to create process directory: %v", err)
	}
	fmt.Println("write pid file: ", pidFile)
	return os.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0644)
}

func TerminateProcess(rootDir, processName string) error {
	pid, err := ReadPIDFile(rootDir, processName)
	if err != nil {
		return fmt.Errorf("failed to read PID file: %v", err)
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process: %v", err)
	}
	return process.Kill()
}

func GetPIDFile(rootDir, processName string) string {
	return filepath.Join(rootDir, "proc", processName+".pid")
}

func ReadPIDFile(rootDir, processName string) (int, error) {
	pidFile := GetPIDFile(rootDir, processName)
	content, err := os.ReadFile(pidFile)
	if err != nil {
		return 0, fmt.Errorf("PID file read error: %v", err)
	}

	pid, err := strconv.Atoi(strings.TrimSpace(string(content)))
	if err != nil {
		return 0, fmt.Errorf("PID conversion error: %v", err)
	}

	return pid, nil
}

func ServiceIsRunning(rootDir, processName string) bool {
	pid, err := ReadPIDFile(rootDir, processName)
	if err != nil {
		return false
	}
	return ProcessIsRunning(pid)
}
