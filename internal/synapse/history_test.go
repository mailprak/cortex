package synapse_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/anoop2811/cortex/internal/synapse"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestHistory(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "History Suite")
}

var _ = Describe("ExecutionHistory", func() {
	var (
		historyManager *synapse.HistoryManager
		testDir        string
		synapseName    string
		err            error
	)

	BeforeEach(func() {
		// Create temporary directory for testing
		testDir, err = os.MkdirTemp("", "cortex-history-test-*")
		Expect(err).NotTo(HaveOccurred())

		synapseName = "test-synapse"
		historyManager = synapse.NewHistoryManager(testDir)
	})

	AfterEach(func() {
		// Clean up test directory
		if testDir != "" {
			os.RemoveAll(testDir)
		}
	})

	Describe("AddExecution", func() {
		Context("when adding a new execution record", func() {
			It("should successfully save the execution record", func() {
				record := synapse.ExecutionRecord{
					ID:          "exec-001",
					SynapseName: synapseName,
					Timestamp:   time.Now(),
					Status:      "success",
					Duration:    time.Second * 5,
					NeuronResults: []synapse.NeuronResult{
						{
							Name:     "neuron-1",
							Status:   "success",
							ExitCode: 0,
							Duration: time.Second * 2,
							Stdout:   "Output from neuron 1",
							Stderr:   "",
						},
					},
				}

				err := historyManager.AddExecution(synapseName, record)
				Expect(err).NotTo(HaveOccurred())

				// Verify file was created
				historyFile := filepath.Join(testDir, synapseName+".json")
				_, err = os.Stat(historyFile)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should append to existing history", func() {
				record1 := synapse.ExecutionRecord{
					ID:          "exec-001",
					SynapseName: synapseName,
					Timestamp:   time.Now(),
					Status:      "success",
					Duration:    time.Second * 3,
				}

				record2 := synapse.ExecutionRecord{
					ID:          "exec-002",
					SynapseName: synapseName,
					Timestamp:   time.Now().Add(time.Minute),
					Status:      "failed",
					Duration:    time.Second * 7,
				}

				err := historyManager.AddExecution(synapseName, record1)
				Expect(err).NotTo(HaveOccurred())

				err = historyManager.AddExecution(synapseName, record2)
				Expect(err).NotTo(HaveOccurred())

				history, err := historyManager.GetHistory(synapseName)
				Expect(err).NotTo(HaveOccurred())
				Expect(history).To(HaveLen(2))
				Expect(history[0].ID).To(Equal("exec-001"))
				Expect(history[1].ID).To(Equal("exec-002"))
			})

			It("should create directory if it doesn't exist", func() {
				nestedDir := filepath.Join(testDir, "nested", "path")
				manager := synapse.NewHistoryManager(nestedDir)

				record := synapse.ExecutionRecord{
					ID:          "exec-001",
					SynapseName: synapseName,
					Timestamp:   time.Now(),
					Status:      "success",
				}

				err := manager.AddExecution(synapseName, record)
				Expect(err).NotTo(HaveOccurred())

				// Verify directory was created
				_, err = os.Stat(nestedDir)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when handling errors", func() {
			It("should handle invalid synapse name", func() {
				record := synapse.ExecutionRecord{
					ID:          "exec-001",
					SynapseName: "",
					Timestamp:   time.Now(),
					Status:      "success",
				}

				err := historyManager.AddExecution("", record)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("synapse name cannot be empty"))
			})
		})
	})

	Describe("GetHistory", func() {
		Context("when retrieving execution history", func() {
			It("should return all execution records", func() {
				// Add multiple records
				for i := 1; i <= 3; i++ {
					record := synapse.ExecutionRecord{
						ID:          "exec-00" + string(rune('0'+i)),
						SynapseName: synapseName,
						Timestamp:   time.Now().Add(time.Duration(i) * time.Minute),
						Status:      "success",
						Duration:    time.Second * time.Duration(i),
					}
					err := historyManager.AddExecution(synapseName, record)
					Expect(err).NotTo(HaveOccurred())
				}

				history, err := historyManager.GetHistory(synapseName)
				Expect(err).NotTo(HaveOccurred())
				Expect(history).To(HaveLen(3))
			})

			It("should return empty slice for non-existent synapse", func() {
				history, err := historyManager.GetHistory("non-existent")
				Expect(err).NotTo(HaveOccurred())
				Expect(history).To(BeEmpty())
			})

			It("should preserve execution order", func() {
				times := []time.Time{
					time.Now(),
					time.Now().Add(time.Minute),
					time.Now().Add(2 * time.Minute),
				}

				for i, t := range times {
					record := synapse.ExecutionRecord{
						ID:          "exec-00" + string(rune('1'+i)),
						SynapseName: synapseName,
						Timestamp:   t,
						Status:      "success",
					}
					err := historyManager.AddExecution(synapseName, record)
					Expect(err).NotTo(HaveOccurred())
				}

				history, err := historyManager.GetHistory(synapseName)
				Expect(err).NotTo(HaveOccurred())
				Expect(history).To(HaveLen(3))

				// Verify chronological order
				for i := range history {
					Expect(history[i].ID).To(Equal("exec-00" + string(rune('1'+i))))
				}
			})
		})
	})

	Describe("GetExecutionLogs", func() {
		Context("when retrieving execution logs", func() {
			It("should return detailed logs for a specific execution", func() {
				record := synapse.ExecutionRecord{
					ID:          "exec-001",
					SynapseName: synapseName,
					Timestamp:   time.Now(),
					Status:      "success",
					Duration:    time.Second * 5,
					NeuronResults: []synapse.NeuronResult{
						{
							Name:     "neuron-1",
							Status:   "success",
							ExitCode: 0,
							Duration: time.Second * 2,
							Stdout:   "Line 1\nLine 2\nLine 3",
							Stderr:   "",
						},
						{
							Name:     "neuron-2",
							Status:   "failed",
							ExitCode: 1,
							Duration: time.Second * 3,
							Stdout:   "Some output",
							Stderr:   "Error occurred",
						},
					},
				}

				err := historyManager.AddExecution(synapseName, record)
				Expect(err).NotTo(HaveOccurred())

				logs, err := historyManager.GetExecutionLogs(synapseName, "exec-001")
				Expect(err).NotTo(HaveOccurred())
				Expect(logs).NotTo(BeNil())
				Expect(logs.ID).To(Equal("exec-001"))
				Expect(logs.NeuronResults).To(HaveLen(2))
				Expect(logs.NeuronResults[0].Stdout).To(ContainSubstring("Line 1"))
				Expect(logs.NeuronResults[1].Stderr).To(ContainSubstring("Error occurred"))
			})

			It("should return error for non-existent execution", func() {
				// First add some history so we can test "execution not found" vs "history not found"
				record := synapse.ExecutionRecord{
					ID:          "exec-001",
					SynapseName: synapseName,
					Timestamp:   time.Now(),
					Status:      "success",
				}
				err := historyManager.AddExecution(synapseName, record)
				Expect(err).NotTo(HaveOccurred())

				// Now look for a non-existent execution
				_, err = historyManager.GetExecutionLogs(synapseName, "non-existent")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("execution not found"))
			})

			It("should return error for non-existent synapse", func() {
				_, err := historyManager.GetExecutionLogs("non-existent", "exec-001")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("history not found"))
			})
		})
	})

	Describe("Concurrent Access", func() {
		Context("when multiple goroutines access history", func() {
			It("should handle concurrent writes safely", func() {
				var wg sync.WaitGroup
				numGoroutines := 10

				for i := 0; i < numGoroutines; i++ {
					wg.Add(1)
					go func(index int) {
						defer wg.Done()
						record := synapse.ExecutionRecord{
							ID:          "exec-" + string(rune('0'+index)),
							SynapseName: synapseName,
							Timestamp:   time.Now(),
							Status:      "success",
							Duration:    time.Second,
						}
						err := historyManager.AddExecution(synapseName, record)
						Expect(err).NotTo(HaveOccurred())
					}(i)
				}

				wg.Wait()

				history, err := historyManager.GetHistory(synapseName)
				Expect(err).NotTo(HaveOccurred())
				Expect(history).To(HaveLen(numGoroutines))
			})

			It("should handle concurrent reads safely", func() {
				// Add some initial records
				for i := 0; i < 5; i++ {
					record := synapse.ExecutionRecord{
						ID:          "exec-00" + string(rune('0'+i)),
						SynapseName: synapseName,
						Timestamp:   time.Now(),
						Status:      "success",
					}
					err := historyManager.AddExecution(synapseName, record)
					Expect(err).NotTo(HaveOccurred())
				}

				var wg sync.WaitGroup
				numGoroutines := 20

				for i := 0; i < numGoroutines; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						history, err := historyManager.GetHistory(synapseName)
						Expect(err).NotTo(HaveOccurred())
						Expect(history).To(HaveLen(5))
					}()
				}

				wg.Wait()
			})

			It("should handle mixed concurrent reads and writes", func() {
				var wg sync.WaitGroup
				numWriters := 5
				numReaders := 10

				// Writers
				for i := 0; i < numWriters; i++ {
					wg.Add(1)
					go func(index int) {
						defer wg.Done()
						record := synapse.ExecutionRecord{
							ID:          "exec-" + string(rune('0'+index)),
							SynapseName: synapseName,
							Timestamp:   time.Now(),
							Status:      "success",
						}
						historyManager.AddExecution(synapseName, record)
					}(i)
				}

				// Readers
				for i := 0; i < numReaders; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						historyManager.GetHistory(synapseName)
					}()
				}

				wg.Wait()

				// Verify final state
				history, err := historyManager.GetHistory(synapseName)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(history)).To(BeNumerically("<=", numWriters))
			})
		})
	})

	Describe("Persistence", func() {
		Context("when data is persisted to disk", func() {
			It("should persist data correctly in JSON format", func() {
				record := synapse.ExecutionRecord{
					ID:          "exec-001",
					SynapseName: synapseName,
					Timestamp:   time.Now(),
					Status:      "success",
					Duration:    time.Second * 5,
					NeuronResults: []synapse.NeuronResult{
						{
							Name:     "neuron-1",
							Status:   "success",
							ExitCode: 0,
							Duration: time.Second * 2,
							Stdout:   "Test output",
							Stderr:   "",
						},
					},
				}

				err := historyManager.AddExecution(synapseName, record)
				Expect(err).NotTo(HaveOccurred())

				// Read file directly and verify JSON format
				historyFile := filepath.Join(testDir, synapseName+".json")
				data, err := os.ReadFile(historyFile)
				Expect(err).NotTo(HaveOccurred())

				var records []synapse.ExecutionRecord
				err = json.Unmarshal(data, &records)
				Expect(err).NotTo(HaveOccurred())
				Expect(records).To(HaveLen(1))
				Expect(records[0].ID).To(Equal("exec-001"))
			})

			It("should survive manager recreation", func() {
				record := synapse.ExecutionRecord{
					ID:          "exec-001",
					SynapseName: synapseName,
					Timestamp:   time.Now(),
					Status:      "success",
				}

				err := historyManager.AddExecution(synapseName, record)
				Expect(err).NotTo(HaveOccurred())

				// Create new manager with same directory
				newManager := synapse.NewHistoryManager(testDir)
				history, err := newManager.GetHistory(synapseName)
				Expect(err).NotTo(HaveOccurred())
				Expect(history).To(HaveLen(1))
				Expect(history[0].ID).To(Equal("exec-001"))
			})
		})
	})

	Describe("Edge Cases", func() {
		Context("when handling special characters in names", func() {
			It("should handle synapse names with special characters", func() {
				specialName := "test-synapse_v1.2.3"
				record := synapse.ExecutionRecord{
					ID:          "exec-001",
					SynapseName: specialName,
					Timestamp:   time.Now(),
					Status:      "success",
				}

				err := historyManager.AddExecution(specialName, record)
				Expect(err).NotTo(HaveOccurred())

				history, err := historyManager.GetHistory(specialName)
				Expect(err).NotTo(HaveOccurred())
				Expect(history).To(HaveLen(1))
			})
		})

		Context("when handling large datasets", func() {
			It("should handle large number of executions", func() {
				numRecords := 100

				for i := 0; i < numRecords; i++ {
					record := synapse.ExecutionRecord{
						ID:          "exec-" + string(rune('0'+i%10)),
						SynapseName: synapseName,
						Timestamp:   time.Now(),
						Status:      "success",
					}
					err := historyManager.AddExecution(synapseName, record)
					Expect(err).NotTo(HaveOccurred())
				}

				history, err := historyManager.GetHistory(synapseName)
				Expect(err).NotTo(HaveOccurred())
				Expect(history).To(HaveLen(numRecords))
			})

			It("should handle large log outputs", func() {
				largeOutput := string(make([]byte, 10000)) // 10KB output

				record := synapse.ExecutionRecord{
					ID:          "exec-001",
					SynapseName: synapseName,
					Timestamp:   time.Now(),
					Status:      "success",
					NeuronResults: []synapse.NeuronResult{
						{
							Name:     "neuron-1",
							Status:   "success",
							ExitCode: 0,
							Stdout:   largeOutput,
						},
					},
				}

				err := historyManager.AddExecution(synapseName, record)
				Expect(err).NotTo(HaveOccurred())

				logs, err := historyManager.GetExecutionLogs(synapseName, "exec-001")
				Expect(err).NotTo(HaveOccurred())
				Expect(len(logs.NeuronResults[0].Stdout)).To(Equal(10000))
			})
		})
	})
})
