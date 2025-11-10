import { useEffect, useRef, useState, useCallback } from 'react';
import { WebSocketMessage, ExecutionLog, ExecutionStatus } from '../types';

interface UseWebSocketOptions {
  url: string;
  onMessage?: (message: WebSocketMessage) => void;
  onConnect?: () => void;
  onDisconnect?: () => void;
  onError?: (error: Event) => void;
  reconnect?: boolean;
  reconnectInterval?: number;
  maxReconnectAttempts?: number;
}

export const useWebSocket = (options: UseWebSocketOptions) => {
  const {
    url,
    onMessage,
    onConnect,
    onDisconnect,
    onError,
    reconnect = true,
    reconnectInterval = 3000,
    maxReconnectAttempts = 5,
  } = options;

  const wsRef = useRef<WebSocket | null>(null);
  const reconnectAttemptsRef = useRef(0);
  const reconnectTimeoutRef = useRef<number>();

  const [isConnected, setIsConnected] = useState(false);
  const [lastMessage, setLastMessage] = useState<WebSocketMessage | null>(null);

  const connect = useCallback(() => {
    try {
      // Clean up existing connection
      if (wsRef.current) {
        wsRef.current.close();
      }

      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      const wsUrl = url.startsWith('ws') ? url : `${protocol}//${window.location.host}${url}`;

      const ws = new WebSocket(wsUrl);
      wsRef.current = ws;

      ws.onopen = () => {
        console.log('WebSocket connected');
        setIsConnected(true);
        reconnectAttemptsRef.current = 0;
        onConnect?.();
      };

      ws.onmessage = (event) => {
        try {
          const message: WebSocketMessage = JSON.parse(event.data);
          setLastMessage(message);
          onMessage?.(message);
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error);
        }
      };

      ws.onerror = (error) => {
        console.error('WebSocket error:', error);
        onError?.(error);
      };

      ws.onclose = () => {
        console.log('WebSocket disconnected');
        setIsConnected(false);
        onDisconnect?.();

        // Attempt reconnection
        if (reconnect && reconnectAttemptsRef.current < maxReconnectAttempts) {
          reconnectAttemptsRef.current += 1;
          console.log(
            `Attempting to reconnect (${reconnectAttemptsRef.current}/${maxReconnectAttempts})...`
          );
          reconnectTimeoutRef.current = setTimeout(() => {
            connect();
          }, reconnectInterval);
        }
      };
    } catch (error) {
      console.error('Failed to create WebSocket connection:', error);
    }
  }, [url, onMessage, onConnect, onDisconnect, onError, reconnect, reconnectInterval, maxReconnectAttempts]);

  const disconnect = useCallback(() => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current);
    }
    if (wsRef.current) {
      wsRef.current.close();
      wsRef.current = null;
    }
    setIsConnected(false);
  }, []);

  const sendMessage = useCallback((message: any) => {
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(message));
    } else {
      console.warn('WebSocket is not connected');
    }
  }, []);

  useEffect(() => {
    connect();

    return () => {
      disconnect();
    };
  }, [connect, disconnect]);

  return {
    isConnected,
    lastMessage,
    sendMessage,
    disconnect,
    reconnect: connect,
  };
};

export const useExecutionLogs = (executionId?: string) => {
  const [logs, setLogs] = useState<ExecutionLog[]>([]);
  const [status, setStatus] = useState<ExecutionStatus | null>(null);

  const handleMessage = useCallback(
    (message: WebSocketMessage) => {
      if (message.type === 'log') {
        const log = message.payload as ExecutionLog;
        if (!executionId || log.neuronId === executionId) {
          setLogs((prev) => [...prev, log]);
        }
      } else if (message.type === 'status') {
        const newStatus = message.payload as ExecutionStatus;
        if (!executionId || newStatus.neuronId === executionId) {
          setStatus(newStatus);
        }
      }
    },
    [executionId]
  );

  const { isConnected, sendMessage } = useWebSocket({
    url: '/ws/executions',
    onMessage: handleMessage,
  });

  const clearLogs = useCallback(() => {
    setLogs([]);
    setStatus(null);
  }, []);

  return {
    logs,
    status,
    isConnected,
    clearLogs,
    sendMessage,
  };
};
