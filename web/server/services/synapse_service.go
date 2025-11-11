package services

import (
	"errors"
	"sync"
	"time"

	"github.com/anoop2811/cortex/web/server/models"
	"github.com/google/uuid"
)

// SynapseService handles synapse business logic and storage
type SynapseService struct {
	synapses map[string]*models.Synapse
	mu       sync.RWMutex
}

// NewSynapseService creates a new synapse service
func NewSynapseService() *SynapseService {
	return &SynapseService{
		synapses: make(map[string]*models.Synapse),
	}
}

// CreateSynapse creates a new synapse
func (s *SynapseService) CreateSynapse(synapse *models.Synapse) (*models.Synapse, error) {
	if synapse.Name == "" {
		return nil, errors.New("synapse name is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Generate ID and timestamps
	synapse.ID = uuid.New().String()
	synapse.CreatedAt = time.Now()
	synapse.UpdatedAt = time.Now()

	// Initialize slices if nil
	if synapse.Nodes == nil {
		synapse.Nodes = []models.SynapseNode{}
	}
	if synapse.Connections == nil {
		synapse.Connections = []models.SynapseConnection{}
	}

	s.synapses[synapse.ID] = synapse
	return synapse, nil
}

// GetSynapse retrieves a synapse by ID
func (s *SynapseService) GetSynapse(id string) (*models.Synapse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	synapse, exists := s.synapses[id]
	if !exists {
		return nil, errors.New("synapse not found")
	}

	return synapse, nil
}

// UpdateSynapse updates an existing synapse
func (s *SynapseService) UpdateSynapse(synapse *models.Synapse) (*models.Synapse, error) {
	if synapse.ID == "" {
		return nil, errors.New("synapse ID is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	existing, exists := s.synapses[synapse.ID]
	if !exists {
		return nil, errors.New("synapse not found")
	}

	// Preserve created time, update modified time
	synapse.CreatedAt = existing.CreatedAt
	synapse.UpdatedAt = time.Now()

	// Initialize slices if nil
	if synapse.Nodes == nil {
		synapse.Nodes = []models.SynapseNode{}
	}
	if synapse.Connections == nil {
		synapse.Connections = []models.SynapseConnection{}
	}

	s.synapses[synapse.ID] = synapse
	return synapse, nil
}

// DeleteSynapse deletes a synapse by ID
func (s *SynapseService) DeleteSynapse(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.synapses[id]; !exists {
		return errors.New("synapse not found")
	}

	delete(s.synapses, id)
	return nil
}

// ListSynapses returns all synapses
func (s *SynapseService) ListSynapses() ([]*models.Synapse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	synapses := make([]*models.Synapse, 0, len(s.synapses))
	for _, synapse := range s.synapses {
		synapses = append(synapses, synapse)
	}

	return synapses, nil
}
