package services

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/anoop2811/cortex/logger"
	"github.com/anoop2811/cortex/web/server/models"
	"gopkg.in/yaml.v2"
)

// NeuronService handles neuron operations
type NeuronService struct {
	logger *logger.StandardLogger
}

// NewNeuronService creates a new NeuronService
func NewNeuronService(log *logger.StandardLogger) *NeuronService {
	return &NeuronService{logger: log}
}

// ListNeurons returns all available neurons
func (s *NeuronService) ListNeurons() ([]models.Neuron, error) {
	neurons := []models.Neuron{}

	// Check current directory and example directory for neurons
	dirs := []string{".", "./example"}

	for _, dir := range dirs {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // Skip errors
			}

			if info.IsDir() && strings.HasSuffix(path, "_neuron") {
				neuronType := "check"
				if strings.Contains(filepath.Base(path), "mutate") {
					neuronType = "mutate"
				}

				neurons = append(neurons, models.Neuron{
					Name:        filepath.Base(path),
					Type:        neuronType,
					Description: fmt.Sprintf("%s neuron", filepath.Base(path)),
					Path:        path,
				})
			}

			return nil
		})

		if err != nil {
			s.logger.Error(err, "Failed to walk directory")
		}
	}

	// Add some example neurons if none found
	if len(neurons) == 0 {
		neurons = []models.Neuron{
			{Name: "example-check", Type: "check", Description: "Example check neuron", Path: "./example/check_neuron"},
			{Name: "example-mutate", Type: "mutate", Description: "Example mutate neuron", Path: "./example/mutate_neuron"},
		}
	}

	return neurons, nil
}

// ListSynapses returns all available synapses
func (s *NeuronService) ListSynapses() ([]models.Synapse, error) {
	synapses := []models.Synapse{}

	// Check current directory and example directory for synapses
	dirs := []string{".", "./example"}

	for _, dir := range dirs {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			if !info.IsDir() && info.Name() == "synapse.yaml" {
				synapse, err := s.loadSynapse(path)
				if err != nil {
					s.logger.Error(err, "Failed to load synapse")
					return nil
				}
				synapses = append(synapses, synapse)
			}

			return nil
		})

		if err != nil {
			s.logger.Error(err, "Failed to walk directory")
		}
	}

	// Add example synapse if none found
	if len(synapses) == 0 {
		synapses = []models.Synapse{
			{
				Name:        "example-workflow",
				Description: "Example workflow synapse",
				Path:        "./example/synapse",
				Neurons:     []string{"check-nginx", "check-redis"},
			},
		}
	}

	return synapses, nil
}

// loadSynapse loads a synapse from a YAML file
func (s *NeuronService) loadSynapse(path string) (models.Synapse, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return models.Synapse{}, err
	}

	var config struct {
		Name       string `yaml:"name"`
		Definition []struct {
			Name string `yaml:"name"`
		} `yaml:"definition"`
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return models.Synapse{}, err
	}

	neurons := []string{}
	for _, def := range config.Definition {
		neurons = append(neurons, def.Name)
	}

	return models.Synapse{
		Name:        config.Name,
		Description: fmt.Sprintf("Synapse: %s", config.Name),
		Path:        filepath.Dir(path),
		Neurons:     neurons,
	}, nil
}
