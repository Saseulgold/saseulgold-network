package kernel

import (
	"os"
	"path/filepath"
	"strconv"
	"syscall"
)

const (
	MASTER         = "master"
	CHAIN_MAKER    = "maker"
	RESOURCE_MINER = "miner"
	COLLECTOR      = "collector"
	PEER_SEARCHER  = "peer_searcher"
	DATA_POOL      = "data_pool"
)

type ProcessManager struct{}

func (pm *ProcessManager) Exists(name string) bool {
	_, err := os.Stat(pm.file(name))
	return !os.IsNotExist(err)
}

func (pm *ProcessManager) file(name string) string {
	return filepath.Join(os.Getenv("DATA_DIR"), name+".pid")
}

func (pm *ProcessManager) Save(pid string) error {
	return os.WriteFile(pm.file(pid), []byte(strconv.Itoa(os.Getpid())), 0644)
}

func (pm *ProcessManager) Delete(pid string) error {
	return os.Remove(pm.file(pid))
}

func (pm *ProcessManager) Pid(name string) (int, error) {
	data, err := os.ReadFile(pm.file(name))
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(data))
}

func (pm *ProcessManager) IsRunning(name string) bool {
	pid, err := pm.Pid(name)
	if err != nil || pid <= 0 {
		return false
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	err = process.Signal(syscall.Signal(0))
	return err == nil
}

func (pm *ProcessManager) Kill(name string) bool {
	if !pm.IsRunning(name) {
		return false
	}

	pid, err := pm.Pid(name)
	if err != nil {
		return false
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	pm.Delete(name)
	process.Kill()
	return true
}
