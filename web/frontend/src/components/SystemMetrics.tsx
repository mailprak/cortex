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
  <div className="group glass rounded-xl shadow-card hover:shadow-card-hover border border-primary-500/20 hover:border-primary-500/40 transition-all duration-300 p-6">
    <div className="flex items-center justify-between mb-4">
      <div className="flex items-center gap-3">
        <div className="relative">
          {icon}
          <div className="absolute inset-0 blur-lg opacity-30 group-hover:opacity-50 transition-opacity"></div>
        </div>
        <span className="font-heading font-semibold text-text-primary">{label}</span>
      </div>
      <span className="text-sm font-bold text-primary-300">{value}</span>
    </div>
    <div className="w-full bg-background-slate/50 rounded-full h-3 overflow-hidden">
      <div
        className={`h-3 rounded-full transition-all duration-500 ${color}`}
        style={{ width: `${Math.min(percentage, 100)}%` }}
      />
    </div>
    <div className="text-xs text-text-secondary mt-2 font-medium">
      {percentage.toFixed(1)}% utilized
    </div>
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
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {[1, 2, 3].map((i) => (
          <div key={i} className="glass rounded-xl shadow-card p-6 animate-pulse">
            <div className="h-6 bg-background-slate/50 rounded-lg mb-4" />
            <div className="h-3 bg-background-slate/50 rounded-full" />
          </div>
        ))}
      </div>
    );
  }

  if (error || !metrics) {
    return (
      <div className="glass border-2 border-red-500/30 rounded-xl p-6">
        <p className="text-red-300 text-lg">{error || 'No metrics available'}</p>
      </div>
    );
  }

  const formatBytes = (bytes: number): string => {
    const gb = bytes / (1024 * 1024 * 1024);
    return `${gb.toFixed(2)} GB`;
  };

  const getColorClass = (percentage: number): string => {
    if (percentage >= 90) return 'bg-gradient-to-r from-red-500 to-red-600';
    if (percentage >= 70) return 'bg-gradient-to-r from-yellow-500 to-orange-500';
    return 'bg-gradient-cyan';
  };

  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
      <MetricCard
        icon={<Cpu className="w-6 h-6 text-accent-blue" />}
        label="CPU Usage"
        value={`${metrics.cpu.cores} cores`}
        percentage={metrics.cpu.usage}
        color={getColorClass(metrics.cpu.usage)}
      />
      <MetricCard
        icon={<MemoryStick className="w-6 h-6 text-primary-400" />}
        label="Memory"
        value={`${formatBytes(metrics.memory.used)} / ${formatBytes(metrics.memory.total)}`}
        percentage={metrics.memory.percentage}
        color={getColorClass(metrics.memory.percentage)}
      />
      <MetricCard
        icon={<HardDrive className="w-6 h-6 text-accent-cyan" />}
        label="Disk Space"
        value={`${formatBytes(metrics.disk.used)} / ${formatBytes(metrics.disk.total)}`}
        percentage={metrics.disk.percentage}
        color={getColorClass(metrics.disk.percentage)}
      />
    </div>
  );
};
