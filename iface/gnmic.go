package iface

import (
	"fmt"
	"log"
	"os/exec"
)

type GnmicIface struct {
	name   string
	config map[string]string
	pids   map[int32]int32
}

func NewGnmicIface() *GnmicIface {
	return &GnmicIface{
		name:   "gnmic",
		config: make(map[string]string),
		pids:   make(map[int32]int32),
	}
}

func (i GnmicIface) GetName() string {
	return i.name
}

func (i GnmicIface) SetConfig(config map[string]string) {
	i.config = config
}

func (i GnmicIface) StartIface() {
	log.Printf("Starting gNMIc interface %s for node %s...", i.name)

	// Prepare the command
	args := []string{"subscribe", "--address", i.name}
	for key, value := range i.config {
		args = append(args, "--"+key, value)
	}
	log.Printf("Starting gNMIc process with args %s...", args)
	proc := exec.Command("gnmic", args...)

	// Start the process
	if err := proc.Start(); err != nil {
		log.Printf("Failed to start gnmic process: %v", err)
		return
	}

	// Store the process ID
	i.pids[int32(proc.Process.Pid)] = int32(proc.Process.Pid)
	log.Printf("gNMIc process started with PID %d", proc.Process.Pid)
}

func (g GnmicIface) StopIface() {
	log.Printf("Stopping gNMIc interface %s...", g.name)

	for pid := range g.pids {
		proc := exec.Command("kill", "-9", fmt.Sprintf("%d", pid))
		if err := proc.Run(); err != nil {
			log.Printf("Failed to stop gnmic process with PID %d: %v", pid, err)
		} else {
			log.Printf("gNMIc process with PID %d stopped", pid)
			delete(g.pids, pid)
		}
	}
}
