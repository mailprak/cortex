import React, { useEffect, useRef } from 'react';
import { Terminal, AlertCircle, Info, AlertTriangle, Bug } from 'lucide-react';
import { useExecutionLogs } from '../hooks/useWebSocket';
import { ExecutionLog } from '../types';

interface ExecutionLogsProps {
  executionId?: string;
  autoScroll?: boolean;
}

const LogLevelIcon: React.FC<{ level: ExecutionLog['level'] }> = ({ level }) => {
  switch (level) {
    case 'error':
      return <AlertCircle className="w-4 h-4 text-red-500" />;
    case 'warn':
      return <AlertTriangle className="w-4 h-4 text-yellow-500" />;
    case 'debug':
      return <Bug className="w-4 h-4 text-purple-500" />;
    default:
      return <Info className="w-4 h-4 text-blue-500" />;
  }
};

const LogEntry: React.FC<{ log: ExecutionLog }> = ({ log }) => {
  const getLevelClass = () => {
    switch (log.level) {
      case 'error':
        return 'border-l-red-500 bg-red-50';
      case 'warn':
        return 'border-l-yellow-500 bg-yellow-50';
      case 'debug':
        return 'border-l-purple-500 bg-purple-50';
      default:
        return 'border-l-blue-500 bg-blue-50';
    }
  };

  return (
    <div className={`flex gap-3 p-2 border-l-4 ${getLevelClass()} font-mono text-sm`}>
      <LogLevelIcon level={log.level} />
      <span className="text-gray-500 min-w-[80px]">
        {new Date(log.timestamp).toLocaleTimeString()}
      </span>
      <span className="flex-1 text-gray-800 break-all">{log.message}</span>
    </div>
  );
};

export const ExecutionLogs: React.FC<ExecutionLogsProps> = ({
  executionId,
  autoScroll = true,
}) => {
  const { logs, status, isConnected, clearLogs } = useExecutionLogs(executionId);
  const logsEndRef = useRef<HTMLDivElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (autoScroll && logsEndRef.current) {
      logsEndRef.current.scrollIntoView({ behavior: 'smooth' });
    }
  }, [logs, autoScroll]);

  const getStatusBadge = () => {
    if (!status) return null;

    const statusColors = {
      running: 'bg-blue-100 text-blue-800',
      completed: 'bg-green-100 text-green-800',
      failed: 'bg-red-100 text-red-800',
    };

    return (
      <div
        data-testid="execution-status"
        className={`px-3 py-1 rounded-full text-xs font-medium ${statusColors[status.status]}`}
      >
        {status.status.toUpperCase()}
        {status.exitCode !== undefined && ` (Exit: ${status.exitCode})`}
      </div>
    );
  };

  return (
    <div className="bg-white rounded-lg shadow-md overflow-hidden">
      {/* Header */}
      <div className="bg-gray-800 text-white px-4 py-3 flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Terminal className="w-5 h-5" />
          <h3 className="font-semibold">Execution Logs</h3>
          {!isConnected && (
            <span className="text-xs bg-red-500 px-2 py-1 rounded">Disconnected</span>
          )}
        </div>
        <div className="flex items-center gap-3">
          {getStatusBadge()}
          <button
            onClick={clearLogs}
            className="text-sm bg-gray-700 hover:bg-gray-600 px-3 py-1 rounded transition-colors"
          >
            Clear
          </button>
        </div>
      </div>

      {/* Logs Container */}
      <div
        ref={containerRef}
        data-testid="log-stream"
        className="h-96 overflow-y-auto bg-gray-50 p-2 space-y-1"
        role="log"
        aria-label="Execution logs"
      >
        {logs.length === 0 ? (
          <div className="flex items-center justify-center h-full text-gray-500">
            <div className="text-center">
              <Terminal className="w-12 h-12 mx-auto mb-2 text-gray-400" />
              <p>No logs yet. Execute a neuron to see logs.</p>
            </div>
          </div>
        ) : (
          <>
            {logs.map((log) => (
              <LogEntry key={log.id} log={log} />
            ))}
            <div ref={logsEndRef} />
          </>
        )}
      </div>

      {/* Footer */}
      {status && (
        <div className="bg-gray-100 px-4 py-2 text-xs text-gray-600 border-t">
          <div className="flex justify-between">
            <span>Started: {new Date(status.startTime).toLocaleString()}</span>
            {status.endTime && (
              <span>
                Duration:{' '}
                {Math.round(
                  (new Date(status.endTime).getTime() - new Date(status.startTime).getTime()) /
                    1000
                )}
                s
              </span>
            )}
          </div>
          {status.error && (
            <div className="mt-1 text-red-600 font-medium">Error: {status.error}</div>
          )}
        </div>
      )}
    </div>
  );
};
