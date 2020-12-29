package synapse

import "github.com/anoop2811/cortex/internal/neuron"

type Synapse struct {
	Name       string              `yaml:"name"`
	Definition []neuron.Definition `yaml:"definition"`
}

type Plan struct {
	Config Config `yaml:"config"`
}

type Steps struct {
	Serial   []string `yaml:"serial"`
	Parallel []string `yaml:"parallel"`
}

type Config struct {
	ExitOnFirstError bool `yaml:"exit_on_first_error"`
}
