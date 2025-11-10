import React, { useState, useCallback } from 'react';
import { DndProvider, useDrag, useDrop } from 'react-dnd';
import { HTML5Backend } from 'react-dnd-html5-backend';
import { Save, Plus, Trash2, Network } from 'lucide-react';
import { Synapse, SynapseNode, Neuron } from '../types';
import { apiClient } from '../api/client';

interface DraggableNeuronProps {
  neuron: Neuron;
}

const DraggableNeuron: React.FC<DraggableNeuronProps> = ({ neuron }) => {
  const [{ isDragging }, drag] = useDrag(() => ({
    type: 'neuron',
    item: { neuron },
    collect: (monitor) => ({
      isDragging: monitor.isDragging(),
    }),
  }));

  return (
    <div
      ref={drag}
      className={`p-3 bg-white border-2 border-gray-200 rounded-lg cursor-move hover:border-primary-500 transition-colors ${
        isDragging ? 'opacity-50' : 'opacity-100'
      }`}
    >
      <div className="font-medium text-sm text-gray-900">{neuron.name}</div>
      <div className="text-xs text-gray-500 mt-1">{neuron.type}</div>
    </div>
  );
};

interface CanvasNodeProps {
  node: SynapseNode;
  onRemove: (id: string) => void;
}

const CanvasNode: React.FC<CanvasNodeProps> = ({ node, onRemove }) => {
  return (
    <div
      className="absolute bg-white border-2 border-primary-500 rounded-lg p-3 shadow-lg"
      style={{ left: node.position.x, top: node.position.y }}
    >
      <div className="flex items-center justify-between gap-2">
        <div>
          <div className="font-medium text-sm">{node.data.label}</div>
          {node.data.description && (
            <div className="text-xs text-gray-500 mt-1">{node.data.description}</div>
          )}
        </div>
        <button
          onClick={() => onRemove(node.id)}
          className="text-red-500 hover:text-red-700"
          aria-label="Remove node"
        >
          <Trash2 className="w-4 h-4" />
        </button>
      </div>
    </div>
  );
};

interface SynapseBuilderProps {
  neurons: Neuron[];
  onSave?: (synapse: Partial<Synapse>) => void;
}

export const SynapseBuilder: React.FC<SynapseBuilderProps> = ({ neurons, onSave }) => {
  const [nodes, setNodes] = useState<SynapseNode[]>([]);
  const [synapseName, setSynapseName] = useState('');
  const [synapseDescription, setSynapseDescription] = useState('');
  const [isSaving, setIsSaving] = useState(false);

  const [{ isOver }, drop] = useDrop(() => ({
    accept: 'neuron',
    drop: (item: { neuron: Neuron }, monitor) => {
      const offset = monitor.getClientOffset();
      if (offset) {
        const canvasRect = document
          .querySelector('[data-testid="synapse-canvas"]')
          ?.getBoundingClientRect();
        if (canvasRect) {
          const newNode: SynapseNode = {
            id: `node-${Date.now()}`,
            type: 'neuron',
            neuronId: item.neuron.id,
            position: {
              x: offset.x - canvasRect.left,
              y: offset.y - canvasRect.top,
            },
            data: {
              label: item.neuron.name,
              description: item.neuron.description,
            },
          };
          setNodes((prev) => [...prev, newNode]);
        }
      }
    },
    collect: (monitor) => ({
      isOver: monitor.isOver(),
    }),
  }));

  const handleRemoveNode = useCallback((id: string) => {
    setNodes((prev) => prev.filter((node) => node.id !== id));
  }, []);

  const handleSave = async () => {
    if (!synapseName.trim()) {
      alert('Please enter a synapse name');
      return;
    }

    const synapse: Partial<Synapse> = {
      name: synapseName,
      description: synapseDescription,
      nodes,
      connections: [], // TODO: Implement connections
    };

    try {
      setIsSaving(true);
      const saved = await apiClient.createSynapse(synapse);
      onSave?.(saved);
      alert('Synapse saved successfully!');
    } catch (error) {
      console.error('Failed to save synapse:', error);
      alert('Failed to save synapse');
    } finally {
      setIsSaving(false);
    }
  };

  const handleClear = () => {
    setNodes([]);
    setSynapseName('');
    setSynapseDescription('');
  };

  return (
    <DndProvider backend={HTML5Backend}>
      <div className="bg-white rounded-lg shadow-md overflow-hidden">
        {/* Header */}
        <div className="bg-gradient-to-r from-primary-500 to-primary-600 text-white px-6 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <Network className="w-6 h-6" />
              <h2 className="text-xl font-bold">Visual Synapse Builder</h2>
            </div>
            <div className="flex gap-2">
              <button
                onClick={handleClear}
                className="flex items-center gap-2 px-4 py-2 bg-white/20 hover:bg-white/30 rounded-lg transition-colors"
              >
                <Trash2 className="w-4 h-4" />
                <span>Clear</span>
              </button>
              <button
                onClick={handleSave}
                disabled={isSaving || nodes.length === 0}
                className="flex items-center gap-2 px-4 py-2 bg-white text-primary-600 hover:bg-gray-100 rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <Save className="w-4 h-4" />
                <span>{isSaving ? 'Saving...' : 'Save'}</span>
              </button>
            </div>
          </div>
        </div>

        <div className="p-6">
          {/* Synapse Info */}
          <div className="mb-6 space-y-3">
            <input
              type="text"
              value={synapseName}
              onChange={(e) => setSynapseName(e.target.value)}
              placeholder="Synapse name..."
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
              aria-label="Synapse name"
            />
            <textarea
              value={synapseDescription}
              onChange={(e) => setSynapseDescription(e.target.value)}
              placeholder="Description (optional)..."
              rows={2}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 resize-none"
              aria-label="Synapse description"
            />
          </div>

          <div className="grid grid-cols-4 gap-6">
            {/* Neuron Palette */}
            <div className="col-span-1">
              <div
                data-testid="neuron-palette"
                className="bg-gray-50 rounded-lg p-4 border-2 border-dashed border-gray-300"
              >
                <h3 className="font-semibold text-gray-700 mb-3 flex items-center gap-2">
                  <Plus className="w-4 h-4" />
                  Neuron Palette
                </h3>
                <div className="space-y-2">
                  {neurons.length === 0 ? (
                    <p className="text-sm text-gray-500 text-center py-4">No neurons available</p>
                  ) : (
                    neurons.map((neuron) => (
                      <DraggableNeuron key={neuron.id} neuron={neuron} />
                    ))
                  )}
                </div>
              </div>
            </div>

            {/* Canvas */}
            <div className="col-span-3">
              <div
                ref={drop}
                data-testid="synapse-canvas"
                className={`relative bg-gray-50 rounded-lg border-2 ${
                  isOver ? 'border-primary-500 bg-primary-50' : 'border-dashed border-gray-300'
                } h-[600px] overflow-hidden transition-colors`}
                role="region"
                aria-label="Synapse canvas"
              >
                {nodes.length === 0 ? (
                  <div className="absolute inset-0 flex items-center justify-center text-gray-400">
                    <div className="text-center">
                      <Network className="w-16 h-16 mx-auto mb-4" />
                      <p className="text-lg font-medium">Drag neurons here to build your synapse</p>
                    </div>
                  </div>
                ) : (
                  nodes.map((node) => (
                    <CanvasNode key={node.id} node={node} onRemove={handleRemoveNode} />
                  ))
                )}
              </div>
              <div className="mt-2 text-sm text-gray-500">
                {nodes.length} node{nodes.length !== 1 ? 's' : ''} added
              </div>
            </div>
          </div>
        </div>
      </div>
    </DndProvider>
  );
};
