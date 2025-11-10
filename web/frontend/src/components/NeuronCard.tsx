import React, { useState } from 'react';
import { Play, Square, Clock, CheckCircle, XCircle, Loader } from 'lucide-react';
import { Neuron, ExecutionStatus } from '../types';
import { apiClient } from '../api/client';

interface NeuronCardProps {
  neuron: Neuron;
  onExecute?: (neuron: Neuron) => void;
  onStatusChange?: (status: ExecutionStatus) => void;
}

export const NeuronCard: React.FC<NeuronCardProps> = ({ neuron, onExecute, onStatusChange }) => {
  const [isExecuting, setIsExecuting] = useState(false);
  const [executionId, setExecutionId] = useState<string | null>(null);

  const getStatusIcon = () => {
    switch (neuron.status) {
      case 'running':
        return <Loader className="w-4 h-4 text-blue-500 animate-spin" />;
      case 'completed':
        return <CheckCircle className="w-4 h-4 text-green-500" />;
      case 'failed':
        return <XCircle className="w-4 h-4 text-red-500" />;
      default:
        return <Clock className="w-4 h-4 text-gray-400" />;
    }
  };

  const getStatusColor = () => {
    switch (neuron.status) {
      case 'running':
        return 'bg-blue-100 text-blue-800';
      case 'completed':
        return 'bg-green-100 text-green-800';
      case 'failed':
        return 'bg-red-100 text-red-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  const handleExecute = async () => {
    try {
      setIsExecuting(true);
      const status = await apiClient.executeNeuron(neuron.id);
      setExecutionId(status.id);
      onExecute?.(neuron);
      onStatusChange?.(status);
    } catch (error) {
      console.error('Failed to execute neuron:', error);
    } finally {
      setIsExecuting(false);
    }
  };

  const handleStop = async () => {
    if (!executionId) return;

    try {
      await apiClient.stopNeuronExecution(neuron.id, executionId);
      setExecutionId(null);
    } catch (error) {
      console.error('Failed to stop execution:', error);
    }
  };

  const isRunning = neuron.status === 'running' || isExecuting;

  return (
    <div
      data-testid="neuron-card"
      className="bg-white rounded-lg shadow-md hover:shadow-lg transition-shadow duration-200 p-6"
    >
      <div className="flex items-start justify-between mb-3">
        <div className="flex-1">
          <h3 className="text-lg font-semibold text-gray-900 mb-1">{neuron.name}</h3>
          <p className="text-sm text-gray-600">{neuron.description}</p>
        </div>
        <div className="flex items-center gap-2">
          {getStatusIcon()}
          <span className={`text-xs px-2 py-1 rounded-full font-medium ${getStatusColor()}`}>
            {neuron.status}
          </span>
        </div>
      </div>

      <div className="flex items-center justify-between mt-4">
        <div className="flex items-center gap-2 text-xs text-gray-500">
          <span className="px-2 py-1 bg-gray-100 rounded">{neuron.type}</span>
          <span>Updated: {new Date(neuron.updatedAt).toLocaleDateString()}</span>
        </div>

        <button
          onClick={isRunning ? handleStop : handleExecute}
          disabled={isExecuting}
          className={`flex items-center gap-2 px-4 py-2 rounded-lg font-medium transition-colors ${
            isRunning
              ? 'bg-red-500 hover:bg-red-600 text-white'
              : 'bg-primary-500 hover:bg-primary-600 text-white disabled:bg-gray-300 disabled:cursor-not-allowed'
          }`}
          aria-label={isRunning ? 'Stop execution' : 'Execute neuron'}
        >
          {isRunning ? (
            <>
              <Square className="w-4 h-4" />
              <span>Stop</span>
            </>
          ) : (
            <>
              <Play className="w-4 h-4" />
              <span>Execute</span>
            </>
          )}
        </button>
      </div>
    </div>
  );
};
