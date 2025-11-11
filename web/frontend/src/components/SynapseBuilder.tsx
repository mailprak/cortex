import React, { useState, useCallback } from 'react';
import { useDrag, useDrop } from 'react-dnd';
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
      className={`p-3 glass border-2 border-primary-500/30 hover:border-primary-500/60 rounded-xl cursor-move transition-all duration-300 ${
        isDragging ? 'opacity-50 scale-95' : 'opacity-100 hover:scale-105'
      }`}
    >
      <div className="font-heading font-semibold text-sm text-text-primary">{neuron.name}</div>
      <div className="text-xs text-text-secondary mt-1">{neuron.type}</div>
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
      className="absolute glass border-2 border-primary-500/40 rounded-xl p-4 shadow-card-hover"
      style={{ left: node.position.x, top: node.position.y }}
    >
      <div className="flex items-center justify-between gap-3">
        <div>
          <div className="font-heading font-semibold text-sm text-text-primary">{node.data.label}</div>
          {node.data.description && (
            <div className="text-xs text-text-secondary mt-1 max-w-[200px] truncate">
              {node.data.description}
            </div>
          )}
        </div>
        <button
          onClick={() => onRemove(node.id)}
          className="text-red-400 hover:text-red-300 transition-colors p-1 hover:bg-red-500/10 rounded"
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
      <div className="glass rounded-2xl shadow-card overflow-hidden border border-primary-500/20">
        {/* Header */}
        <div className="relative bg-gradient-purple">
          <div className="relative px-6 py-5">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <Network className="w-7 h-7 text-white" />
                <h2 className="text-2xl font-heading font-bold text-white">Visual Synapse Builder</h2>
              </div>
              <div className="flex gap-3">
                <button
                  onClick={handleClear}
                  className="flex items-center gap-2 px-5 py-2.5 bg-white/20 hover:bg-white/30 backdrop-blur-sm rounded-pill transition-all duration-300 text-white font-medium"
                >
                  <Trash2 className="w-4 h-4" />
                  <span>Clear</span>
                </button>
                <button
                  onClick={handleSave}
                  disabled={isSaving || nodes.length === 0}
                  className="flex items-center gap-2 px-5 py-2.5 bg-white text-primary-600 hover:bg-gray-100 rounded-pill transition-all duration-300 disabled:opacity-50 disabled:cursor-not-allowed font-medium shadow-lg"
                >
                  <Save className="w-4 h-4" />
                  <span>{isSaving ? 'Saving...' : 'Save'}</span>
                </button>
              </div>
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
              className="w-full px-5 py-3 bg-background-card border-2 border-primary-500/30 focus:border-primary-500 rounded-xl text-text-primary placeholder-text-muted focus:outline-none transition-colors"
              aria-label="Synapse name"
            />
            <textarea
              value={synapseDescription}
              onChange={(e) => setSynapseDescription(e.target.value)}
              placeholder="Description (optional)..."
              rows={2}
              className="w-full px-5 py-3 bg-background-card border-2 border-primary-500/30 focus:border-primary-500 rounded-xl text-text-primary placeholder-text-muted focus:outline-none transition-colors resize-none"
              aria-label="Synapse description"
            />
          </div>

          <div className="grid grid-cols-4 gap-6">
            {/* Neuron Palette */}
            <div className="col-span-1">
              <div
                data-testid="neuron-palette"
                className="bg-background-card rounded-xl p-5 border-2 border-dashed border-primary-500/30"
              >
                <h3 className="font-heading font-bold text-text-primary mb-4 flex items-center gap-2">
                  <Plus className="w-5 h-5 text-accent-cyan" />
                  Neuron Palette
                </h3>
                <div className="space-y-3">
                  {neurons.length === 0 ? (
                    <p className="text-sm text-text-muted text-center py-8">No neurons available</p>
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
                className={`relative bg-background-card rounded-xl border-2 ${
                  isOver ? 'border-primary-500 bg-primary-500/5' : 'border-dashed border-primary-500/30'
                } h-[600px] overflow-hidden transition-all duration-300`}
                role="region"
                aria-label="Synapse canvas"
              >
                {nodes.length === 0 ? (
                  <div className="absolute inset-0 flex items-center justify-center text-text-muted">
                    <div className="text-center">
                      <Network className="w-20 h-20 mx-auto mb-6 opacity-30" />
                      <p className="text-xl font-heading font-semibold text-text-secondary">
                        Drag neurons here to build your synapse
                      </p>
                    </div>
                  </div>
                ) : (
                  nodes.map((node) => (
                    <CanvasNode key={node.id} node={node} onRemove={handleRemoveNode} />
                  ))
                )}
              </div>
              <div className="mt-3 text-sm text-text-secondary font-medium">
                {nodes.length} node{nodes.length !== 1 ? 's' : ''} added
              </div>
            </div>
          </div>
        </div>
      </div>
  );
};
