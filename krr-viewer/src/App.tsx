import { useState } from 'react';
import './App.css';
import { Tabs } from './components/Tabs';
import { DnsView } from './views/DnsView';
import { KrrView } from './views/KrrView';
import { NodeOptimizerView } from './views/NodeOptimizerView';
import { PlatformOptimizerView } from './views/PlatformOptimizerView';
import { RegistryView } from './views/RegistryView';
import { StorageView } from './views/StorageView';
import { VpcView } from './views/VpcView';

const TABS = ['Pod Resources', 'VPC', 'Storage', 'Registry', 'DNS', 'Node', 'Platform'];

function App() {
  const [activeTab, setActiveTab] = useState(TABS[0]);

  const renderActiveTab = () => {
    switch (activeTab) {
      case 'Pod Resources':
        return <KrrView />;
      case 'VPC':
        return <VpcView />;
      case 'Storage':
        return <StorageView />;
      case 'Registry':
        return <RegistryView />;
      case 'DNS':
        return <DnsView />;
      case 'Node':
        return <NodeOptimizerView />;
      case 'Platform':
        return <PlatformOptimizerView />;
      default:
        return <KrrView />;
    }
  }

  return (
    <div className="app-container">
      <div className="sticky-header">
        <h1 className="main-header">Kubeconomist</h1>
        <Tabs
          tabs={TABS}
          activeTab={activeTab}
          onTabClick={setActiveTab}
        />
      </div>
      <main className="main-content">
        {renderActiveTab()}
      </main>
    </div>
  )
}

export default App;
