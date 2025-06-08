import React from 'react';
import './Tabs.css';

interface TabProps {
  label: string;
  isActive: boolean;
  onClick: () => void;
}

const Tab: React.FC<TabProps> = ({ label, isActive, onClick }) => (
  <button
    className={`tab ${isActive ? 'active' : ''}`}
    onClick={onClick}
  >
    {label}
  </button>
);

interface TabsProps {
  activeTab: string;
  onTabClick: (label: string) => void;
  tabs: string[];
}

export const Tabs: React.FC<TabsProps> = ({ activeTab, onTabClick, tabs }) => {
  return (
    <div className="tabs-container">
      {tabs.map(tab => (
        <Tab
          key={tab}
          label={tab}
          isActive={activeTab === tab}
          onClick={() => onTabClick(tab)}
        />
      ))}
    </div>
  );
}; 
