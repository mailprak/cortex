package services_test

import (
	"testing"

	"github.com/anoop2811/cortex/web/server/models"
	"github.com/anoop2811/cortex/web/server/services"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSynapseService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Synapse Service Suite")
}

var _ = Describe("SynapseService", func() {
	var service *services.SynapseService

	BeforeEach(func() {
		service = services.NewSynapseService()
	})

	Describe("CreateSynapse", func() {
		It("should create a synapse with a generated ID", func() {
			synapse := &models.Synapse{
				Name:        "test-synapse",
				Description: "A test synapse",
				Nodes:       []models.SynapseNode{},
				Connections: []models.SynapseConnection{},
			}

			created, err := service.CreateSynapse(synapse)
			Expect(err).NotTo(HaveOccurred())
			Expect(created.ID).NotTo(BeEmpty())
			Expect(created.Name).To(Equal("test-synapse"))
			Expect(created.CreatedAt).NotTo(BeZero())
			Expect(created.UpdatedAt).NotTo(BeZero())
		})

		It("should return error for empty name", func() {
			synapse := &models.Synapse{
				Name: "",
			}

			_, err := service.CreateSynapse(synapse)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("name"))
		})

		It("should store multiple synapses", func() {
			synapse1 := &models.Synapse{Name: "synapse1", Nodes: []models.SynapseNode{}, Connections: []models.SynapseConnection{}}
			synapse2 := &models.Synapse{Name: "synapse2", Nodes: []models.SynapseNode{}, Connections: []models.SynapseConnection{}}

			created1, _ := service.CreateSynapse(synapse1)
			created2, _ := service.CreateSynapse(synapse2)

			Expect(created1.ID).NotTo(Equal(created2.ID))

			synapses, _ := service.ListSynapses()
			Expect(synapses).To(HaveLen(2))
		})
	})

	Describe("GetSynapse", func() {
		It("should retrieve a synapse by ID", func() {
			synapse := &models.Synapse{
				Name:        "get-test",
				Nodes:       []models.SynapseNode{},
				Connections: []models.SynapseConnection{},
			}

			created, _ := service.CreateSynapse(synapse)
			retrieved, err := service.GetSynapse(created.ID)

			Expect(err).NotTo(HaveOccurred())
			Expect(retrieved.ID).To(Equal(created.ID))
			Expect(retrieved.Name).To(Equal("get-test"))
		})

		It("should return error for non-existent ID", func() {
			_, err := service.GetSynapse("non-existent-id")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not found"))
		})
	})

	Describe("UpdateSynapse", func() {
		It("should update an existing synapse", func() {
			synapse := &models.Synapse{
				Name:        "update-test",
				Description: "Original",
				Nodes:       []models.SynapseNode{},
				Connections: []models.SynapseConnection{},
			}

			created, _ := service.CreateSynapse(synapse)

			updated := &models.Synapse{
				ID:          created.ID,
				Name:        "update-test",
				Description: "Updated",
				Nodes: []models.SynapseNode{
					{ID: "node-1", Type: "neuron", NeuronID: "test"},
				},
				Connections: []models.SynapseConnection{},
			}

			result, err := service.UpdateSynapse(updated)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Description).To(Equal("Updated"))
			Expect(result.Nodes).To(HaveLen(1))
			Expect(result.UpdatedAt).To(BeTemporally(">", created.UpdatedAt))
		})

		It("should return error for non-existent synapse", func() {
			synapse := &models.Synapse{
				ID:          "non-existent",
				Name:        "test",
				Nodes:       []models.SynapseNode{},
				Connections: []models.SynapseConnection{},
			}

			_, err := service.UpdateSynapse(synapse)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("DeleteSynapse", func() {
		It("should delete an existing synapse", func() {
			synapse := &models.Synapse{
				Name:        "delete-test",
				Nodes:       []models.SynapseNode{},
				Connections: []models.SynapseConnection{},
			}

			created, _ := service.CreateSynapse(synapse)
			err := service.DeleteSynapse(created.ID)

			Expect(err).NotTo(HaveOccurred())

			_, getErr := service.GetSynapse(created.ID)
			Expect(getErr).To(HaveOccurred())
		})

		It("should return error for non-existent ID", func() {
			err := service.DeleteSynapse("non-existent")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("ListSynapses", func() {
		It("should return empty list initially", func() {
			synapses, err := service.ListSynapses()
			Expect(err).NotTo(HaveOccurred())
			Expect(synapses).To(BeEmpty())
		})

		It("should return all created synapses", func() {
			for i := 0; i < 3; i++ {
				synapse := &models.Synapse{
					Name:        "list-test",
					Nodes:       []models.SynapseNode{},
					Connections: []models.SynapseConnection{},
				}
				service.CreateSynapse(synapse)
			}

			synapses, err := service.ListSynapses()
			Expect(err).NotTo(HaveOccurred())
			Expect(synapses).To(HaveLen(3))
		})
	})
})
