package types

import (
	_ "embed"
	"github.com/go-yaml/yaml"
	"github.com/stretchr/testify/assert"
	"testing"
)

//go:embed testdata/simple_topo.yaml
var expectedTopology string

func TestTopology(t *testing.T) {
	builder := NewTopologyBuilder()
	builder.Name("srlceos01")
	builder.AddNode("nokia_srlinux", "srlinux")
	builder.AddNode("arista_ceos", "ceos")
	builder.AddLink("nokia_srlinux", "arista_ceos", "e1-1", "eth1")

	topo := builder.Build()
	output, err := yaml.Marshal(topo)
	if err != nil {
		t.Fatalf("Failed to marshal topology: %s", err)
	}

	t.Log(string(output))
	assert.Equal(t, expectedTopology, string(output), "Topology is not as expected")
}
