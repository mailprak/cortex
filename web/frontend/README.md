# Cortex Web UI - Frontend

Modern React-based web interface for Cortex neural network orchestration system.

## Features

- **Dashboard**: Real-time neuron library with execution controls
- **System Metrics**: Live CPU, Memory, and Disk monitoring
- **Execution Logs**: WebSocket-powered real-time log streaming
- **Visual Synapse Builder**: Drag-and-drop interface for building neural workflows
- **Responsive Design**: Mobile-first with hamburger menu
- **Accessibility**: ARIA labels, keyboard navigation, semantic HTML

## Tech Stack

- **React 18** with TypeScript
- **Vite** for fast development and optimized builds
- **TailwindCSS** for styling
- **React DnD** for drag-and-drop functionality
- **React Router** for navigation
- **Axios** for API communication
- **WebSocket** for real-time updates
- **Lucide React** for icons

## Getting Started

### Prerequisites

- Node.js 18+
- npm or yarn

### Installation

```bash
npm install
```

### Development

```bash
npm run dev
```

The app will be available at `http://localhost:3000`

### Build

```bash
npm run build
```

### Preview Production Build

```bash
npm run preview
```

### Type Checking

```bash
npm run typecheck
```

### Linting

```bash
npm run lint
```

## Project Structure

```
src/
├── components/       # React components
│   ├── Dashboard.tsx
│   ├── NeuronCard.tsx
│   ├── ExecutionLogs.tsx
│   ├── SynapseBuilder.tsx
│   └── SystemMetrics.tsx
├── hooks/           # Custom React hooks
│   └── useWebSocket.ts
├── api/             # API client
│   └── client.ts
├── types/           # TypeScript definitions
│   └── index.ts
├── styles/          # Global styles
│   └── index.css
├── App.tsx          # Main application component
└── main.tsx         # Entry point
```

## API Configuration

The frontend expects the backend API at `http://localhost:8080`. Configure via Vite proxy in `vite.config.ts`.

## Accessibility Features

- Semantic HTML elements
- ARIA labels and roles
- Keyboard navigation support
- Focus management
- Screen reader friendly

## Performance

- Target load time: <2 seconds
- Code splitting with React.lazy
- Optimized bundle size
- Efficient re-renders with React.memo

## Testing

E2E tests are located in `/tests/e2e/` and use Playwright.

Test selectors used:
- `data-testid="neuron-card"` - Neuron cards
- `data-testid="log-stream"` - Execution logs container
- `data-testid="execution-status"` - Status display
- `data-testid="neuron-palette"` - Synapse builder palette
- `data-testid="synapse-canvas"` - Synapse builder canvas

## License

MIT
