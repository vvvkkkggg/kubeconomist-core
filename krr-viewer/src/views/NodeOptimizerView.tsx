import { unparse } from 'papaparse';
import React, { useMemo, useState } from 'react';
import { SortableHeader } from '../components/SortableHeader';
import { ViewHeader } from '../components/ViewHeader';
import { mockNodeOptimizerRecommendations } from '../data/node-optimizer-mock-data';
import { useSort } from '../hooks/useSort';
import type { NodeOptimizerRecommendation } from '../types';

export const NodeOptimizerView: React.FC = () => {
    const [searchQuery, setSearchQuery] = useState('');

    const optimizableRecs = useMemo(() => {
        return mockNodeOptimizerRecommendations.filter((rec: NodeOptimizerRecommendation) => rec.currentPrice > rec.desiredPrice);
    }, []);

    const filteredRecs = useMemo(() => {
        const lowerCaseQuery = searchQuery.toLowerCase();
        return optimizableRecs.filter((rec: NodeOptimizerRecommendation) =>
            rec.instanceId.toLowerCase().includes(lowerCaseQuery) ||
            rec.cloudId.toLowerCase().includes(lowerCaseQuery) ||
            rec.folderId.toLowerCase().includes(lowerCaseQuery)
        );
    }, [searchQuery, optimizableRecs]);

    const { items: sortedRecs, requestSort, sortKey, sortDirection } = useSort(filteredRecs, 'instanceId', 'ascending');

    const sortProps = {
        currentSortKey: sortKey,
        direction: sortDirection,
        onRequestSort: requestSort,
    };

    const handleExport = () => {
        const dataToExport = sortedRecs.map((rec: NodeOptimizerRecommendation) => ({
            'Instance ID': rec.instanceId,
            'Cloud ID': rec.cloudId,
            'Folder ID': rec.folderId,
            'Current Cores': rec.currentCores,
            'Desired Cores': rec.desiredCores,
            'Memory (GB)': rec.currentMemoryGB,
            'Current Monthly Cost ($)': rec.currentPrice,
            'Desired Monthly Cost ($)': rec.desiredPrice,
            'Monthly Savings ($)': rec.currentPrice - rec.desiredPrice,
        }));
        const csv = unparse(dataToExport);
        const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' });
        const link = document.createElement('a');
        link.href = URL.createObjectURL(blob);
        link.setAttribute('download', 'node-optimizer-recommendations.csv');
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
    };

    return (
        <div>
            <h2 className="view-header">Node Optimizer Recommendations ({sortedRecs.length})</h2>
            <ViewHeader
                onSearch={setSearchQuery}
                onHideEmptyToggle={() => {}}
                showHideEmpty={false}
                onExport={handleExport}
            />
            <div className="table-container">
                <table>
                    <thead>
                        <tr>
                            <SortableHeader sortKey="instanceId" {...sortProps}>Instance ID</SortableHeader>
                            <SortableHeader sortKey="cloudId" {...sortProps}>Cloud ID</SortableHeader>
                            <SortableHeader sortKey="folderId" {...sortProps}>Folder ID</SortableHeader>
                            <SortableHeader sortKey="currentCores" {...sortProps}>Current Cores</SortableHeader>
                            <SortableHeader sortKey="desiredCores" {...sortProps}>Desired Cores</SortableHeader>
                            <SortableHeader sortKey="currentMemoryGB" {...sortProps}>Memory (GB)</SortableHeader>
                            <SortableHeader sortKey="currentPrice" {...sortProps}>Current Cost ($)</SortableHeader>
                            <SortableHeader sortKey="desiredPrice" {...sortProps}>Desired Cost ($)</SortableHeader>
                            <SortableHeader sortKey="savings" {...sortProps}>Savings ($)</SortableHeader>
                        </tr>
                    </thead>
                    <tbody>
                        {sortedRecs.map((rec: NodeOptimizerRecommendation) => {
                            const savings = rec.currentPrice - rec.desiredPrice;
                            return (
                                <tr key={rec.id}>
                                    <td>{rec.instanceId}</td>
                                    <td>{rec.cloudId}</td>
                                    <td>{rec.folderId}</td>
                                    <td>{rec.currentCores}</td>
                                    <td>{rec.desiredCores}</td>
                                    <td>{rec.currentMemoryGB}</td>
                                    <td>{rec.currentPrice.toFixed(2)}</td>
                                    <td>{rec.desiredPrice.toFixed(2)}</td>
                                    <td>{savings.toFixed(2)}</td>
                                </tr>
                            );
                        })}
                    </tbody>
                </table>
            </div>
        </div>
    );
}; 
