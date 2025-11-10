export interface Neuron {
  id: string;
  name: string;
  description: string;
  type: string;
  status: 'idle' | 'running' | 'completed' | 'failed';
  createdAt: string;
  updatedAt: string;
  metadata?: Record<string, any>;
}

export interface ExecutionLog {
  id: string;
  neuronId: string;
  timestamp: string;
  level: 'info' | 'warn' | 'error' | 'debug';
  message: string;
}

export interface ExecutionStatus {
  id: string;
  neuronId: string;
  status: 'running' | 'completed' | 'failed';
  startTime: string;
  endTime?: string;
  exitCode?: number;
  error?: string;
}

export interface SystemMetrics {
  cpu: {
    usage: number;
    cores: number;
  };
  memory: {
    used: number;
    total: number;
    percentage: number;
  };
  disk: {
    used: number;
    total: number;
    percentage: number;
  };
  uptime: number;
}

export interface SynapseNode {
  id: string;
  type: 'neuron' | 'input' | 'output';
  neuronId?: string;
  position: {
    x: number;
    y: number;
  };
  data: {
    label: string;
    description?: string;
  };
}

export interface SynapseConnection {
  id: string;
  source: string;
  target: string;
  type: 'data' | 'control';
}

export interface Synapse {
  id: string;
  name: string;
  description: string;
  nodes: SynapseNode[];
  connections: SynapseConnection[];
  createdAt: string;
  updatedAt: string;
}

export interface WebSocketMessage {
  type: 'log' | 'status' | 'metrics';
  payload: ExecutionLog | ExecutionStatus | SystemMetrics;
}
