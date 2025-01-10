package builder

import (
	"bytes"
	"fmt"
	"github.com/Lachstec/digsinet-ng/types"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"os/exec"
)

type ClabBuilder struct {
}

func NewClabBuilder() *ClabBuilder {
	return &ClabBuilder{}
}

func (b *ClabBuilder) DeployTopology(topology types.Topology) error {

	log.Info().
		Str("Builder", b.Id()).
		Msg("Deploying Topology...")

	// Prepare the command
	proc := exec.Command("clab", "deploy", "--topo", "-")

	// Connect stdin for the process
	stdin, err := proc.StdinPipe()
	if err != nil {
		log.Error().
			Str("Builder", b.Id()).
			Msg("Failed to create stdin pipe")
		return fmt.Errorf("failed to get stdin pipe: %w", err)
	}

	// Store stdout of the process
	stdout := bytes.NewBuffer([]byte{})
	proc.Stdout = stdout

	// Store stderr of the process
	stderr := bytes.NewBuffer([]byte{})
	proc.Stderr = stderr

	// Start the process
	if err = proc.Start(); err != nil {
		log.Error().
			Str("Builder", b.Id()).
			Msg("Failed to create stdin pipe")
		return fmt.Errorf("failed to start process: %w", err)
	}

	// Serialize the topology spec
	topologySpec, err := yaml.Marshal(topology)
	if err != nil {
		log.Error().
			Str("Builder", b.Id()).
			Msg("Failed to marshal topology spec")
		return fmt.Errorf("failed to marshal topology: %w", err)
	}

	log.Debug().
		Str("Builder", b.Id()).
		Msg("Topology Spec to be sent: ")

	// Write the topology spec to stdin
	_, err = stdin.Write(topologySpec)
	if err != nil {
		log.Error().
			Str("Builder", b.Id()).
			Msg("Failed to write topology spec to stdin")
		return fmt.Errorf("failed to write to stdin: %w", err)
	}

	// Close stdin to signal the end of input
	if err := stdin.Close(); err != nil {
		log.Error().
			Str("Builder", b.Id()).
			Msg("Failed to close stdin")
		return fmt.Errorf("failed to close stdin: %w", err)
	}

	// Wait for the process to complete
	if err = proc.Wait(); err != nil {
		log.Error().
			Str("Builder", b.Id()).
			Msg("Failed to deploy topology")
		return fmt.Errorf("process finished with error: %w stdout: %s stderr: %s", err, stdout, stderr)
	}

	log.Info().
		Str("Builder", b.Id()).
		Msg("Successfully deployed topology")

	// Start interfaces of the nodes
	//log.Print("Starting interfaces of the nodes...")
	//for _, node := range topology.Nodes {
	//	for _, iface := range node.Ifaces {
	//		log.Printf("Starting iface %s of node %s...", iface.GetName(), node.Name)
	//		iface.StartIface()
	//		if err := proc.Run(); err != nil {
	//			return fmt.Errorf("failed to start iface %s of node %s: %w stdout: %s stderr: %s", iface.GetName(), node.Name, err, stdout, stderr)
	//		}
	//	}
	//}
	//
	//log.Print("Interface configuration completed successfully.")

	return nil
}

func (b *ClabBuilder) DestroyTopology(topology types.Topology) error {
	log.Info().
		Str("Builder", b.Id()).
		Msg("Destroying Topology...")

	// Prepare the command
	proc := exec.Command("clab", "destroy", "--topo", "-")

	// Connect stdin for the process
	stdin, err := proc.StdinPipe()
	if err != nil {
		log.Error().
			Str("Builder", b.Id()).
			Msg("Failed to create stdin pipe")
		return fmt.Errorf("failed to get stdin pipe: %w", err)
	}

	// Store stdout of the process
	stdout := bytes.NewBuffer([]byte{})
	proc.Stdout = stdout

	// Store stderr of the process
	stderr := bytes.NewBuffer([]byte{})
	proc.Stderr = stderr

	// Start the process
	if err = proc.Start(); err != nil {
		log.Error().
			Str("Builder", b.Id()).
			Msg("Failed to start process")
		return fmt.Errorf("failed to start process: %w", err)
	}

	// Serialize the topology spec
	topologySpec, err := yaml.Marshal(topology)
	if err != nil {
		log.Error().
			Str("Builder", b.Id()).
			Msg("Failed to marshal topology spec")
		return fmt.Errorf("failed to marshal topology: %w", err)
	}

	// Debug: print the topology spec to ensure correctness
	log.Debug().
		Str("Builder", b.Id()).
		Msg("Topology Spec to be sent: ")

	// Write the topology spec to stdin
	_, err = stdin.Write(topologySpec)
	if err != nil {
		log.Error().
			Str("Builder", b.Id()).
			Msg("Failed to write topology spec to stdin")
		return fmt.Errorf("failed to write to stdin: %w", err)
	}

	// Close stdin to signal the end of input
	if err := stdin.Close(); err != nil {
		log.Error().
			Str("Builder", b.Id()).
			Msg("Failed to close stdin")
		return fmt.Errorf("failed to close stdin: %w", err)
	}

	// Wait for the process to complete
	if err = proc.Wait(); err != nil {
		log.Error().
			Str("Builder", b.Id()).
			Msg("Failed to deploy topology")
		return fmt.Errorf("process finished with error: %w stdout: %s stderr: %s", err, stdout, stderr)
	}

	log.Info().
		Str("Builder", b.Id()).
		Msg("Successfully deployed topology")
	return nil
}

func (b *ClabBuilder) Id() string {
	return "Containerlab"
}
