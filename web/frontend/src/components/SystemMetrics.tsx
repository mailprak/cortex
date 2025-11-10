import React, { useEffect, useState } from 'react';
import { Cpu, HardDrive, MemoryStick } from 'lucide-react';
import { SystemMetrics as SystemMetricsType } from '../types';
import { apiClient } from '../api/client';

const MetricCard: React.FC<{
  icon: React.ReactNode;
  label: string;
  value: string;
  percentage: number;
  color: string;
}> = ({ icon, label, value, percentage, color }) => (
  <div className="bg-white rounded-lg shadow-md p-4">
    <div className="flex items-center justify-between mb-2">
      <div className="flex items-center gap-2">
        {icon}
        <span className="font-medium text-gray-700">{label}</span>
      </div>
      <span className="text-sm font-semibold text-gray-900">{value}</span>
    </div>
    <div className="w-full bg-gray-200 rounded-full h-2">
      <div
        className={`h-2 rounded-full transition-all duration-300 ${color}`}
        style={{ width: `${Math.min(percentage, 100)}%` }}
      />
    </div>
    <div className="text-xs text-gray-500 mt-1">{percentage.toFixed(1)}%</div>
  </div>
);

export const SystemMetrics: React.FC = () => {
  const [metrics, setMetrics] = useState<SystemMetricsType | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchMetrics = async () => {
      try {
        const data = await apiClient.getSystemMetrics();
        setMetrics(data);
        setError(null);
      } catch (err) {
        setError('Failed to fetch system metrics');
        console.error('Error fetching metrics:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchMetrics();
    const interval = setInterval(fetchMetrics, 5000); // Update every 5 seconds

    return () => clearInterval(interval);
  }, []);

  if (loading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {[1, 2, 3].map((i) => (
          <div key={i} className="bg-white rounded-lg shadow-md p-4 animate-pulse">
            <div className="h-6 bg-gray-200 rounded mb-2" />
            <div className="h-2 bg-gray-200 rounded" />
          </div>
        ))}
      </div>
    );
  }

  if (error || !metrics) {
    return (
      <div className="bg-red-50 border border-red-200 rounded-lg p-4">
        <p className="text-red-700">{error || 'No metrics available'}</p>
      </div>
    );
  }

  const formatBytes = (bytes: number): string => {
    const gb = bytes / (1024 * 1024 * 1024);
    return `${gb.toFixed(2)} GB`;
  };

  const getColorClass = (percentage: number): string => {
    if (percentage >= 90) return 'bg-red-500';
    if (percentage >= 70) return 'bg-yellow-500';
    return 'bg-green-500';
  };

  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
      <MetricCard
        icon={<Cpu className="w-5 h-5 text-blue-500" />}
        label="CPU Usage"
        value={`${metrics.cpu.cores} cores`}
        percentage={metrics.cpu.usage}
        color={getColorClass(metrics.cpu.usage)}
      />
      <MetricCard
        icon={<MemoryStick className="w-5 h-5 text-purple-500" />}
        label="Memory"
        value={`${formatBytes(metrics.memory.used)} / ${formatBytes(metrics.memory.total)}`}
        percentage={metrics.memory.percentage}
        color={getColorClass(metrics.memory.percentage)}
      />
      <MetricCard
        icon={<HardDrive className="w-5 h-5 text-green-500" />}
        label="Disk Space"
        value={`${formatBytes(metrics.disk.used)} / ${formatBytes(metrics.disk.total)}`}
        percentage={metrics.disk.percentage}
        color={getColorClass(metrics.disk.percentage)}
      />
    </div>
  );
};
