import React from 'react';
import './SubHeaderTabs.css';

interface SubTabProps {
  label: string;
  isActive: boolean;
  onClick: () => void;
}

const SubTab: React.FC<SubTabProps> = ({ label, isActive, onClick }) => (
  <button
    className={`sub-tab ${isActive ? 'active' : ''}`}
    onClick={onClick}
  >
    {label}
  </button>
);

interface SubHeaderTabsProps {
  activeTab: string;
  onTabClick: (label: string) => void;
  tabs: string[];
}

export const SubHeaderTabs: React.FC<SubHeaderTabsProps> = ({ activeTab, onTabClick, tabs }) => {
  return (
    <div className="sub-tabs-container">
      {tabs.map(tab => (
        <SubTab
          key={tab}
          label={tab}
          isActive={activeTab === tab}
          onClick={() => onTabClick(tab)}
        />
      ))}
    </div>
  );
}; 
