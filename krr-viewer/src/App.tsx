import { useState } from 'react';
import './App.css';
import { Tabs } from './components/Tabs';
import { KrrView } from './views/KrrView';
import { RegistryView } from './views/RegistryView';
import { StorageView } from './views/StorageView';
import { VpcView } from './views/VpcView';

const TABS = ['KRR', 'VPC', 'Storage', 'Registry'];

function App() {
  const [activeTab, setActiveTab] = useState(TABS[0]);

  const renderActiveTab = () => {
    switch (activeTab) {
      case 'KRR':
        return <KrrView />;
      case 'VPC':
        return <VpcView />;
      case 'Storage':
        return <StorageView />;
      case 'Registry':
        return <RegistryView />;
      default:
        return <KrrView />;
    }
  }

  return (
    <div className="app-container">
      <h1 className="main-header">Kubeconomist</h1>
      <Tabs
        tabs={TABS}
        activeTab={activeTab}
        onTabClick={setActiveTab}
      />
      <main>
        {renderActiveTab()}
      </main>
    </div>
  )
}

export default App;
