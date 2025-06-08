import React, { useState } from 'react';
import { SortableHeader } from '../components/SortableHeader';
import { SubHeaderTabs } from '../components/SubHeaderTabs';
import { mockSubnetRecommendations } from '../data/subnet-mock-data';
import { mockVpcRecommendations } from '../data/vpc-mock-data';
import { useSort } from '../hooks/useSort';
import type { SubnetRecommendation, VpcRecommendation } from '../types';

const TABS = ['Unused IPs', 'Subnet Utilization'];

const UnusedIpsTable: React.FC<{ items: VpcRecommendation[] }> = ({ items }) => {
    const { items: sortedRecs, requestSort, sortKey, sortDirection } = useSort(items, 'severity', 'descending');
    const sortProps = { currentSortKey: sortKey, direction: sortDirection, onRequestSort: requestSort };

    return (
        <div className="table-container">
            <table>
                <thead>
                    <tr>
                        <SortableHeader sortKey="ipAddress" {...sortProps}>IP Address</SortableHeader>
                        <SortableHeader sortKey="resourceName" {...sortProps}>Resource Name</SortableHeader>
                        <SortableHeader sortKey="region" {...sortProps}>Region</SortableHeader>
                        <SortableHeader sortKey="unusedForDays" {...sortProps}>Unused For (Days)</SortableHeader>
                        <SortableHeader sortKey="severity" {...sortProps}>Severity</SortableHeader>
                    </tr>
                </thead>
                <tbody>
                    {sortedRecs.map((rec) => (
                        <tr key={rec.id}>
                            <td>{rec.ipAddress}</td>
                            <td>{rec.resourceName}</td>
                            <td>{rec.region}</td>
                            <td>{rec.unusedForDays}</td>
                            <td><span className={`severity-pill severity-${rec.severity.toLowerCase()}`}>{rec.severity}</span></td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
}

const SubnetUtilizationTable: React.FC<{ items: SubnetRecommendation[] }> = ({ items }) => {
    const { items: sortedRecs, requestSort, sortKey, sortDirection } = useSort(items, 'utilization', 'descending');
    const sortProps = { currentSortKey: sortKey, direction: sortDirection, onRequestSort: requestSort };

    return (
        <div className="table-container">
            <table>
                <thead>
                    <tr>
                        <SortableHeader sortKey="subnetId" {...sortProps}>Subnet ID</SortableHeader>
                        <SortableHeader sortKey="cidrBlock" {...sortProps}>CIDR Block</SortableHeader>
                        <SortableHeader sortKey="region" {...sortProps}>Region</SortableHeader>
                        <SortableHeader sortKey="utilization" {...sortProps}>Utilization (%)</SortableHeader>
                        <SortableHeader sortKey="severity" {...sortProps}>Severity</SortableHeader>
                    </tr>
                </thead>
                <tbody>
                    {sortedRecs.map((rec) => (
                        <tr key={rec.id}>
                            <td>{rec.subnetId}</td>
                            <td>{rec.cidrBlock}</td>
                            <td>{rec.region}</td>
                            <td>{rec.utilization.toFixed(1)}</td>
                            <td><span className={`severity-pill severity-${rec.severity.toLowerCase()}`}>{rec.severity}</span></td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
}

export const VpcView: React.FC = () => {
    const [activeTab, setActiveTab] = useState(TABS[0]);
    const count = activeTab === TABS[0] ? mockVpcRecommendations.length : mockSubnetRecommendations.length;

    return (
        <div>
            <h2 className="view-header">VPC Recommendations ({count})</h2>
            <SubHeaderTabs tabs={TABS} activeTab={activeTab} onTabClick={setActiveTab} />
            {activeTab === TABS[0] && <UnusedIpsTable items={mockVpcRecommendations} />}
            {activeTab === TABS[1] && <SubnetUtilizationTable items={mockSubnetRecommendations} />}
        </div>
    );
}; 
