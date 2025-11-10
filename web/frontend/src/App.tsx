import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, Link, Navigate } from 'react-router-dom';
import { Menu, X, LayoutDashboard, Network, Settings as SettingsIcon } from 'lucide-react';
import { Dashboard } from './components/Dashboard';
import { SynapseBuilder } from './components/SynapseBuilder';
import { apiClient } from './api/client';
import { Neuron } from './types';

const Navigation: React.FC<{ isMobileMenuOpen: boolean; onToggleMobileMenu: () => void }> = ({
  isMobileMenuOpen,
  onToggleMobileMenu,
}) => {
  const navItems = [
    { path: '/', label: 'Dashboard', icon: LayoutDashboard },
    { path: '/synapse-builder', label: 'Synapse Builder', icon: Network },
    { path: '/settings', label: 'Settings', icon: SettingsIcon },
  ];

  return (
    <>
      {/* Desktop Navigation */}
      <nav
        className="hidden md:flex items-center gap-6"
        aria-label="Main navigation"
      >
        {navItems.map((item) => (
          <Link
            key={item.path}
            to={item.path}
            className="flex items-center gap-2 text-white hover:text-gray-200 transition-colors"
          >
            <item.icon className="w-5 h-5" />
            <span>{item.label}</span>
          </Link>
        ))}
      </nav>

      {/* Mobile Menu Button */}
      <button
        onClick={onToggleMobileMenu}
        className="md:hidden text-white"
        aria-label="Toggle mobile menu"
      >
        {isMobileMenuOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
      </button>

      {/* Mobile Navigation */}
      {isMobileMenuOpen && (
        <div className="md:hidden absolute top-16 left-0 right-0 bg-primary-600 shadow-lg z-50">
          <nav className="flex flex-col p-4 gap-2" aria-label="Mobile navigation">
            {navItems.map((item) => (
              <Link
                key={item.path}
                to={item.path}
                onClick={onToggleMobileMenu}
                className="flex items-center gap-3 text-white hover:bg-primary-700 px-4 py-3 rounded-lg transition-colors"
              >
                <item.icon className="w-5 h-5" />
                <span>{item.label}</span>
              </Link>
            ))}
          </nav>
        </div>
      )}
    </>
  );
};

const Settings: React.FC = () => {
  return (
    <div className="min-h-screen bg-gray-100">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <h1 className="text-2xl font-bold text-gray-900 mb-6">Settings</h1>
        <div className="bg-white rounded-lg shadow-md p-6">
          <p className="text-gray-600">Settings page coming soon...</p>
        </div>
      </div>
    </div>
  );
};

const SynapseBuilderPage: React.FC = () => {
  const [neurons, setNeurons] = useState<Neuron[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchNeurons = async () => {
      try {
        const data = await apiClient.getNeurons();
        setNeurons(data);
      } catch (error) {
        console.error('Failed to fetch neurons:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchNeurons();
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center h-screen">
        <p className="text-gray-600">Loading...</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-100">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <SynapseBuilder neurons={neurons} />
      </div>
    </div>
  );
};

const App: React.FC = () => {
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const [loadStartTime] = useState(Date.now());

  useEffect(() => {
    // Track load time for performance monitoring
    const loadTime = Date.now() - loadStartTime;
    console.log(`App loaded in ${loadTime}ms`);

    // Ensure load time is under 2 seconds as per requirements
    if (loadTime > 2000) {
      console.warn('App load time exceeded 2 seconds');
    }
  }, [loadStartTime]);

  const toggleMobileMenu = () => {
    setIsMobileMenuOpen((prev) => !prev);
  };

  return (
    <Router>
      <div className="min-h-screen bg-gray-100">
        {/* Header */}
        <header className="bg-gradient-to-r from-primary-600 to-primary-700 shadow-lg">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="flex items-center justify-between h-16">
              <div className="flex items-center">
                <Link to="/" className="flex items-center gap-2">
                  <Network className="w-8 h-8 text-white" />
                  <h1 className="text-xl font-bold text-white">Cortex</h1>
                </Link>
              </div>
              <Navigation
                isMobileMenuOpen={isMobileMenuOpen}
                onToggleMobileMenu={toggleMobileMenu}
              />
            </div>
          </div>
        </header>

        {/* Routes */}
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/synapse-builder" element={<SynapseBuilderPage />} />
          <Route path="/settings" element={<Settings />} />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </div>
    </Router>
  );
};

export default App;
