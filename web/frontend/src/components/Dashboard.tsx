import React, { useEffect, useState } from 'react';
import { Loader, AlertCircle, RefreshCw } from 'lucide-react';
import { Neuron, ExecutionStatus } from '../types';
import { apiClient } from '../api/client';
import { NeuronCard } from './NeuronCard';
import { SystemMetrics } from './SystemMetrics';
import { ExecutionLogs } from './ExecutionLogs';

export const Dashboard: React.FC = () => {
  const [neurons, setNeurons] = useState<Neuron[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedExecution, setSelectedExecution] = useState<string | null>(null);
  const [refreshing, setRefreshing] = useState(false);

  const fetchNeurons = async () => {
    try {
      setRefreshing(true);
      const data = await apiClient.getNeurons();
      setNeurons(data);
      setError(null);
    } catch (err) {
      setError('Failed to load neurons');
      console.error('Error fetching neurons:', err);
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  };

  useEffect(() => {
    fetchNeurons();

    // Set up periodic refresh
    const interval = setInterval(fetchNeurons, 10000); // Refresh every 10 seconds

    return () => clearInterval(interval);
  }, []);

  const handleNeuronExecute = (neuron: Neuron) => {
    // Update neuron status optimistically
    setNeurons((prev) =>
      prev.map((n) => (n.id === neuron.id ? { ...n, status: 'running' } : n))
    );
  };

  const handleStatusChange = (status: ExecutionStatus) => {
    setSelectedExecution(status.id);

    // Update neuron status when execution completes
    if (status.status !== 'running') {
      setNeurons((prev) =>
        prev.map((n) => (n.id === status.neuronId ? { ...n, status: status.status } : n))
      );
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-screen">
        <div className="text-center">
          <Loader className="w-12 h-12 animate-spin mx-auto mb-4 text-primary-500" />
          <p className="text-gray-600">Loading dashboard...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center h-screen">
        <div className="bg-red-50 border border-red-200 rounded-lg p-6 max-w-md">
          <AlertCircle className="w-12 h-12 text-red-500 mx-auto mb-4" />
          <h3 className="text-lg font-semibold text-red-900 mb-2">Error Loading Dashboard</h3>
          <p className="text-red-700 mb-4">{error}</p>
          <button
            onClick={fetchNeurons}
            className="w-full bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded-lg transition-colors"
          >
            Retry
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-100">
      {/* Header */}
      <header className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <div className="flex items-center justify-between">
            <h1 className="text-2xl font-bold text-gray-900">Cortex Dashboard</h1>
            <button
              onClick={fetchNeurons}
              disabled={refreshing}
              className="flex items-center gap-2 px-4 py-2 bg-primary-500 hover:bg-primary-600 text-white rounded-lg transition-colors disabled:opacity-50"
              aria-label="Refresh dashboard"
            >
              <RefreshCw className={`w-4 h-4 ${refreshing ? 'animate-spin' : ''}`} />
              <span>Refresh</span>
            </button>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* System Metrics */}
        <section className="mb-8">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">System Metrics</h2>
          <SystemMetrics />
        </section>

        {/* Neuron Library */}
        <section className="mb-8">
          <h2
            className="text-lg font-semibold text-gray-900 mb-4"
            aria-label="Neuron library"
          >
            Neuron Library
          </h2>
          {neurons.length === 0 ? (
            <div className="bg-white rounded-lg shadow-md p-12 text-center">
              <p className="text-gray-500 mb-4">No neurons available</p>
              <button className="bg-primary-500 hover:bg-primary-600 text-white px-6 py-2 rounded-lg transition-colors">
                Create Your First Neuron
              </button>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {neurons.map((neuron) => (
                <NeuronCard
                  key={neuron.id}
                  neuron={neuron}
                  onExecute={handleNeuronExecute}
                  onStatusChange={handleStatusChange}
                />
              ))}
            </div>
          )}
        </section>

        {/* Execution Logs */}
        {selectedExecution && (
          <section>
            <h2 className="text-lg font-semibold text-gray-900 mb-4">Live Execution Logs</h2>
            <ExecutionLogs executionId={selectedExecution} />
          </section>
        )}
      </main>
    </div>
  );
};
