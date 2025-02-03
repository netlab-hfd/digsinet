package builder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os/exec"

	"github.com/Lachstec/digsinet-ng/iface"
	"github.com/Lachstec/digsinet-ng/types"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
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

// needs also iface type as argument
func (b *ClabBuilder) StartNodeIface(topology types.Topology, node string) (string, error) {
	log.Info().
		Str("Builder", b.Id()).
		Msg("Starting Node Interface...")

	// Prepare the command
	proc := exec.Command("clab", "inspect", "--name", topology.Name, "--format", "json")

	// Store stdout of the process
	stdout := bytes.NewBuffer([]byte{})
	proc.Stdout = stdout

	// Store stderr of the process
	stderr := bytes.NewBuffer([]byte{})
	proc.Stderr = stderr

	// Start the process
	if err := proc.Start(); err != nil {
		log.Error().
			Str("Builder", b.Id()).
			Msg("Failed to start node iface process")
		return "", fmt.Errorf("failed to start node iface process: %w", err)
	}

	// Wait for the process to complete
	if err := proc.Wait(); err != nil {
		log.Error().
			Str("Builder", b.Id()).
			Msg("Failed to wait for node iface process")
		return "", fmt.Errorf("process finished with error: %w stdout: %s stderr: %s", err, stdout, stderr)
	}

	// Get node address from stdout
	var containers map[string]interface{}
	if err := json.Unmarshal(stdout.Bytes(), &containers); err != nil {
		log.Error().
			Str("Builder", b.Id()).
			Msg("Failed to unmarshal inspect result")
		return "", fmt.Errorf("failed to unmarshal inspect result: %w", err)
	}

	for _, container := range containers["containers"].([]interface{}) {
		containerData, ok := container.(map[string]interface{})
		if !ok {
			continue
		}

		name, ok := containerData["name"].(string)
		if !ok {
			continue
		}

		if name == "clab-"+topology.Name+"-"+node {
			ipv4AddressString, ok := containerData["ipv4_address"].(string)
			if !ok {
				log.Error().
					Str("Builder", b.Id()).
					Msg("Failed to get IP address of the node")
				return "", fmt.Errorf("failed to get IP address of the node")
			}
			ipv4Address, _, err := net.ParseCIDR(ipv4AddressString)
			if err != nil {
				log.Error().
					Str("Builder", b.Id()).
					Msg("Failed to parse IP address of the node")
				return "", fmt.Errorf("failed to parse IP address of the node: %w", err)
			}
			log.Info().
				Str("Builder", b.Id()).
				Str("Node", node).
				Str("IP", ipv4Address.String()).
				Msg("Successfully got IP address of the node")

			// run gNMI path subscription
			gh, err := iface.NewGNMIHandler()
			if err != nil {
				log.Error().
					Str("Builder", b.Id()).
					Msg("Failed to create GNMI handler")
				return "", fmt.Errorf("failed to create GNMI handler: %w", err)
			}
			subscriptionID, err := gh.SubscribeAndPublish(ipv4Address.String(), []string{"interfaces"}, topology.Name+"-"+node)
			if err != nil {
				log.Error().
					Str("Builder", b.Id()).
					Msg("Failed to subscribe and publish")
				return "", fmt.Errorf("failed to subscribe and publish: %w", err)
			}
			return subscriptionID, nil
		}
	}

	log.Info().
		Str("Builder", b.Id()).
		Msg("Successfully added iface to node topology")
	return "", nil
}

func (b *ClabBuilder) StopNodeIface(topology types.Topology, node string, subscriptionID string) error {
	log.Info().
		Str("Builder", b.Id()).
		Msg("Stopping Node Interface...")

	// run gNMI path subscription
	gh, err := iface.NewGNMIHandler()
	if err != nil {
		log.Error().
			Str("Builder", b.Id()).
			Msg("Failed to create GNMI handler")
		return fmt.Errorf("failed to create GNMI handler: %w", err)
	}
	err = gh.Unsubscribe(topology.Name+"-"+node, subscriptionID)
	if err != nil {
		log.Error().
			Str("Builder", b.Id()).
			Msg("Failed to unsubscribe")
		return fmt.Errorf("failed to unsubscribe: %w", err)
	}

	log.Info().
		Str("Builder", b.Id()).
		Msg("Successfully removed iface from node topology")
	return nil
}

func (b *ClabBuilder) Id() string {
	return "Containerlab"
}
