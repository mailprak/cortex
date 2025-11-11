package acceptance_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anoop2811/cortex/logger"
	"github.com/anoop2811/cortex/web/server"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSynapseBuilder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Synapse Builder API Suite")
}

var _ = Describe("Synapse Builder API", Label("acceptance", "web-api", "synapse-builder"), func() {
	var (
		srv        *server.Server
		testServer *httptest.Server
		apiURL     string
	)

	BeforeEach(func() {
		log := logger.NewLogger(0) // Error level only
		srv = server.NewServer("localhost", 0, log)
		testServer = httptest.NewServer(srv.Router())
		apiURL = testServer.URL
	})

	AfterEach(func() {
		testServer.Close()
	})

	Describe("Creating a new synapse", func() {
		Context("when creating a synapse with valid data", func() {
			It("should create and return the synapse with an ID", func() {
				// RED: This test will fail because the endpoint doesn't exist yet
				synapseData := map[string]interface{}{
					"name":        "test-synapse",
					"description": "A test synapse for health checks",
					"nodes": []map[string]interface{}{
						{
							"id":       "node-1",
							"type":     "neuron",
							"neuronId": "check-nginx",
							"position": map[string]int{"x": 100, "y": 100},
							"data": map[string]string{
								"label":       "Check Nginx",
								"description": "Verify Nginx is running",
							},
						},
						{
							"id":       "node-2",
							"type":     "neuron",
							"neuronId": "check-database",
							"position": map[string]int{"x": 300, "y": 100},
							"data": map[string]string{
								"label":       "Check Database",
								"description": "Verify database connectivity",
							},
						},
					},
					"connections": []map[string]interface{}{
						{
							"id":     "conn-1",
							"source": "node-1",
							"target": "node-2",
						},
					},
				}

				body, _ := json.Marshal(synapseData)
				resp, err := http.Post(apiURL+"/api/synapses", "application/json", bytes.NewBuffer(body))
				Expect(err).NotTo(HaveOccurred())
				defer resp.Body.Close()

				Expect(resp.StatusCode).To(Equal(http.StatusCreated))

				var result map[string]interface{}
				bodyBytes, _ := io.ReadAll(resp.Body)
				json.Unmarshal(bodyBytes, &result)

				Expect(result["id"]).NotTo(BeEmpty())
				Expect(result["name"]).To(Equal("test-synapse"))
				Expect(result["description"]).To(Equal("A test synapse for health checks"))
				Expect(result["nodes"]).To(HaveLen(2))
				Expect(result["connections"]).To(HaveLen(1))
			})

			It("should save the synapse and make it available via GET", func() {
				// Create synapse
				synapseData := map[string]interface{}{
					"name":        "persistent-synapse",
					"description": "Testing persistence",
					"nodes": []map[string]interface{}{
						{
							"id":       "node-1",
							"type":     "neuron",
							"neuronId": "test-neuron",
							"position": map[string]int{"x": 50, "y": 50},
							"data":     map[string]string{"label": "Test Node"},
						},
					},
					"connections": []map[string]interface{}{},
				}

				body, _ := json.Marshal(synapseData)
				createResp, err := http.Post(apiURL+"/api/synapses", "application/json", bytes.NewBuffer(body))
				Expect(err).NotTo(HaveOccurred())
				defer createResp.Body.Close()

				var created map[string]interface{}
				bodyBytes, _ := io.ReadAll(createResp.Body)
				json.Unmarshal(bodyBytes, &created)
				synapseID := created["id"].(string)

				// Retrieve synapse
				getResp, err := http.Get(apiURL + "/api/synapses/" + synapseID)
				Expect(err).NotTo(HaveOccurred())
				defer getResp.Body.Close()

				Expect(getResp.StatusCode).To(Equal(http.StatusOK))

				var retrieved map[string]interface{}
				bodyBytes, _ = io.ReadAll(getResp.Body)
				json.Unmarshal(bodyBytes, &retrieved)

				Expect(retrieved["id"]).To(Equal(synapseID))
				Expect(retrieved["name"]).To(Equal("persistent-synapse"))
				Expect(retrieved["nodes"]).To(HaveLen(1))
			})
		})

		Context("when creating a synapse with invalid data", func() {
			It("should return 400 for missing name", func() {
				synapseData := map[string]interface{}{
					"description": "Missing name",
					"nodes":       []map[string]interface{}{},
				}

				body, _ := json.Marshal(synapseData)
				resp, err := http.Post(apiURL+"/api/synapses", "application/json", bytes.NewBuffer(body))
				Expect(err).NotTo(HaveOccurred())
				defer resp.Body.Close()

				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				var result map[string]interface{}
				bodyBytes, _ := io.ReadAll(resp.Body)
				json.Unmarshal(bodyBytes, &result)

				Expect(result["error"]).To(ContainSubstring("name"))
			})

			It("should return 400 for invalid JSON", func() {
				invalidJSON := bytes.NewBufferString("{invalid json}")
				resp, err := http.Post(apiURL+"/api/synapses", "application/json", invalidJSON)
				Expect(err).NotTo(HaveOccurred())
				defer resp.Body.Close()

				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
			})
		})
	})

	Describe("Updating a synapse", func() {
		It("should update an existing synapse", func() {
			// Create initial synapse
			synapseData := map[string]interface{}{
				"name":        "update-test",
				"description": "Original description",
				"nodes":       []map[string]interface{}{},
				"connections": []map[string]interface{}{},
			}

			body, _ := json.Marshal(synapseData)
			createResp, _ := http.Post(apiURL+"/api/synapses", "application/json", bytes.NewBuffer(body))
			defer createResp.Body.Close()

			var created map[string]interface{}
			bodyBytes, _ := io.ReadAll(createResp.Body)
			json.Unmarshal(bodyBytes, &created)
			synapseID := created["id"].(string)

			// Update synapse
			updateData := map[string]interface{}{
				"name":        "update-test",
				"description": "Updated description",
				"nodes": []map[string]interface{}{
					{
						"id":       "new-node",
						"type":     "neuron",
						"neuronId": "added-neuron",
						"position": map[string]int{"x": 200, "y": 200},
						"data":     map[string]string{"label": "Added Node"},
					},
				},
				"connections": []map[string]interface{}{},
			}

			updateBody, _ := json.Marshal(updateData)
			req, _ := http.NewRequest(http.MethodPut, apiURL+"/api/synapses/"+synapseID, bytes.NewBuffer(updateBody))
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			updateResp, err := client.Do(req)
			Expect(err).NotTo(HaveOccurred())
			defer updateResp.Body.Close()

			Expect(updateResp.StatusCode).To(Equal(http.StatusOK))

			var updated map[string]interface{}
			bodyBytes, _ = io.ReadAll(updateResp.Body)
			json.Unmarshal(bodyBytes, &updated)

			Expect(updated["description"]).To(Equal("Updated description"))
			Expect(updated["nodes"]).To(HaveLen(1))
		})
	})

	Describe("Deleting a synapse", func() {
		It("should delete an existing synapse", func() {
			// Create synapse
			synapseData := map[string]interface{}{
				"name":        "delete-test",
				"description": "To be deleted",
				"nodes":       []map[string]interface{}{},
				"connections": []map[string]interface{}{},
			}

			body, _ := json.Marshal(synapseData)
			createResp, _ := http.Post(apiURL+"/api/synapses", "application/json", bytes.NewBuffer(body))
			defer createResp.Body.Close()

			var created map[string]interface{}
			bodyBytes, _ := io.ReadAll(createResp.Body)
			json.Unmarshal(bodyBytes, &created)
			synapseID := created["id"].(string)

			// Delete synapse
			req, _ := http.NewRequest(http.MethodDelete, apiURL+"/api/synapses/"+synapseID, nil)
			client := &http.Client{}
			deleteResp, err := client.Do(req)
			Expect(err).NotTo(HaveOccurred())
			defer deleteResp.Body.Close()

			Expect(deleteResp.StatusCode).To(Equal(http.StatusNoContent))

			// Verify deletion
			getResp, _ := http.Get(apiURL + "/api/synapses/" + synapseID)
			defer getResp.Body.Close()

			Expect(getResp.StatusCode).To(Equal(http.StatusNotFound))
		})
	})

	Describe("Listing all synapses", func() {
		It("should return all created synapses", func() {
			// Create multiple synapses
			for i := 1; i <= 3; i++ {
				synapseData := map[string]interface{}{
					"name":        "list-test-" + string(rune('0'+i)),
					"description": "Test synapse " + string(rune('0'+i)),
					"nodes":       []map[string]interface{}{},
					"connections": []map[string]interface{}{},
				}

				body, _ := json.Marshal(synapseData)
				resp, _ := http.Post(apiURL+"/api/synapses", "application/json", bytes.NewBuffer(body))
				resp.Body.Close()
			}

			// List all synapses
			listResp, err := http.Get(apiURL + "/api/synapses")
			Expect(err).NotTo(HaveOccurred())
			defer listResp.Body.Close()

			Expect(listResp.StatusCode).To(Equal(http.StatusOK))

			var synapses []map[string]interface{}
			bodyBytes, _ := io.ReadAll(listResp.Body)
			json.Unmarshal(bodyBytes, &synapses)

			Expect(len(synapses)).To(BeNumerically(">=", 3))
		})
	})
})
