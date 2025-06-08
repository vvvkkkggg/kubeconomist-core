import { unparse } from 'papaparse';
import React, { useMemo, useState } from 'react';
import { SortableHeader } from '../components/SortableHeader';
import { ViewHeader } from '../components/ViewHeader';
import { mockPlatformOptimizerRecommendations } from '../data/platform-optimizer-mock-data';
import { useSort } from '../hooks/useSort';
import type { PlatformOptimizerRecommendation } from '../types';

const FolderLink: React.FC<{ folderId: string }> = ({ folderId }) => (
  <a href={`https://console.yandex.cloud/folders/${folderId}`} target="_blank" rel="noopener noreferrer">
    {folderId}
  </a>
);

export const PlatformOptimizerView: React.FC = () => {
    const [searchQuery, setSearchQuery] = useState('');

    const optimizableRecs = useMemo(() => {
        return mockPlatformOptimizerRecommendations.filter((rec: PlatformOptimizerRecommendation) => rec.savings > 0);
    }, []);

    const filteredRecs = useMemo(() => {
        const lowerCaseQuery = searchQuery.toLowerCase();
        return optimizableRecs.filter((rec: PlatformOptimizerRecommendation) =>
            rec.nodeGroupId.toLowerCase().includes(lowerCaseQuery) ||
            rec.cloudId.toLowerCase().includes(lowerCaseQuery) ||
            rec.folderId.toLowerCase().includes(lowerCaseQuery) ||
            rec.currentPlatform.toLowerCase().includes(lowerCaseQuery) ||
            rec.desiredPlatform.toLowerCase().includes(lowerCaseQuery)
        );
    }, [searchQuery, optimizableRecs]);

    const { items: sortedRecs, requestSort, sortKey, sortDirection } = useSort<PlatformOptimizerRecommendation>(filteredRecs, 'savings', 'descending');

    const sortProps = {
        currentSortKey: sortKey,
        direction: sortDirection,
        onRequestSort: requestSort,
    };

    const handleExport = () => {
        const dataToExport = sortedRecs.map((rec: PlatformOptimizerRecommendation) => ({
            'Node Group ID': rec.nodeGroupId,
            'Cloud ID': rec.cloudId,
            'Folder ID': rec.folderId,
            'Current Platform': rec.currentPlatform,
            'Desired Platform': rec.desiredPlatform,
            'Current Monthly Cost ($)': rec.currentPrice,
            'Desired Monthly Cost ($)': rec.desiredPrice,
            'Monthly Savings ($)': rec.savings,
        }));
        const csv = unparse(dataToExport);
        const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' });
        const link = document.createElement('a');
        link.href = URL.createObjectURL(blob);
        link.setAttribute('download', 'platform-optimizer-recommendations.csv');
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
    };

    return (
        <div>
            <h2 className="view-header">Platform Optimizer Recommendations ({sortedRecs.length})</h2>
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
                            <SortableHeader<PlatformOptimizerRecommendation> sortKey="nodeGroupId" {...sortProps}>Node Group ID</SortableHeader>
                            <SortableHeader<PlatformOptimizerRecommendation> sortKey="cloudId" {...sortProps}>Cloud ID</SortableHeader>
                            <SortableHeader<PlatformOptimizerRecommendation> sortKey="folderId" {...sortProps}>Folder ID</SortableHeader>
                            <SortableHeader<PlatformOptimizerRecommendation> sortKey="currentPlatform" {...sortProps}>Current Platform</SortableHeader>
                            <SortableHeader<PlatformOptimizerRecommendation> sortKey="desiredPlatform" {...sortProps}>Desired Platform</SortableHeader>
                            <SortableHeader<PlatformOptimizerRecommendation> sortKey="currentPrice" {...sortProps}>Current Cost ($)</SortableHeader>
                            <SortableHeader<PlatformOptimizerRecommendation> sortKey="desiredPrice" {...sortProps}>Desired Cost ($)</SortableHeader>
                            <SortableHeader<PlatformOptimizerRecommendation> sortKey="savings" {...sortProps}>Savings ($)</SortableHeader>
                        </tr>
                    </thead>
                    <tbody>
                        {sortedRecs.map((rec: PlatformOptimizerRecommendation) => (
                            <tr key={rec.id}>
                                <td>{rec.nodeGroupId}</td>
                                <td>{rec.cloudId}</td>
                                <td><FolderLink folderId={rec.folderId} /></td>
                                <td>{rec.currentPlatform}</td>
                                <td>{rec.desiredPlatform}</td>
                                <td>{rec.currentPrice.toFixed(2)}</td>
                                <td>{rec.desiredPrice.toFixed(2)}</td>
                                <td>{rec.savings.toFixed(2)}</td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
}; 
