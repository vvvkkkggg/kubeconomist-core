import { unparse } from 'papaparse';
import React, { useMemo, useState } from 'react';
import { SortableHeader } from '../components/SortableHeader';
import { ViewHeader } from '../components/ViewHeader';
import { mockRegistryRecommendations } from '../data/registry-mock-data';
import { useSort } from '../hooks/useSort';

const timeSince = (date: string): string => {
    const seconds = Math.floor((new Date().getTime() - new Date(date).getTime()) / 1000);
    let interval = seconds / 31536000;
    if (interval > 1) return Math.floor(interval) + " years ago";
    interval = seconds / 2592000;
    if (interval > 1) return Math.floor(interval) + " months ago";
    interval = seconds / 86400;
    if (interval > 1) return Math.floor(interval) + " days ago";
    return "Today";
};

export const RegistryView: React.FC = () => {
    const [searchQuery, setSearchQuery] = useState('');

    const filteredRecs = useMemo(() => {
        const lowerCaseQuery = searchQuery.toLowerCase();
        return mockRegistryRecommendations.filter(rec =>
            rec.imageName.toLowerCase().includes(lowerCaseQuery) ||
            rec.tags.some(tag => tag.toLowerCase().includes(lowerCaseQuery))
        );
    }, [searchQuery]);

    const { items: sortedRecs, requestSort, sortKey, sortDirection } = useSort(filteredRecs, 'severity', 'descending');

    const sortProps = {
        currentSortKey: sortKey,
        direction: sortDirection,
        onRequestSort: requestSort,
    };

    const handleExport = () => {
        const dataToExport = sortedRecs.map(rec => ({
            ImageName: rec.imageName,
            Tags: rec.tags.join(', '),
            SizeMB: rec.sizeMB,
            LastUsed: rec.lastUsed,
            Severity: rec.severity,
        }));
        const csv = unparse(dataToExport);
        const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' });
        const link = document.createElement('a');
        link.href = URL.createObjectURL(blob);
        link.setAttribute('download', 'registry-recommendations.csv');
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
    };

    return (
        <div>
            <h2 className="view-header">Container Registry Recommendations ({sortedRecs.length})</h2>
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
                            <SortableHeader sortKey="imageName" {...sortProps}>Image Name</SortableHeader>
                            <th>Tags</th>
                            <SortableHeader sortKey="sizeMB" {...sortProps}>Size (MB)</SortableHeader>
                            <SortableHeader sortKey="lastUsed" {...sortProps}>Last Used</SortableHeader>
                            <SortableHeader sortKey="severity" {...sortProps}>Severity</SortableHeader>
                        </tr>
                    </thead>
                    <tbody>
                        {sortedRecs.map((rec) => (
                            <tr key={rec.id}>
                                <td>{rec.imageName}</td>
                                <td>{rec.tags.join(', ')}</td>
                                <td>{rec.sizeMB}</td>
                                <td>{timeSince(rec.lastUsed)}</td>
                                <td><span className={`severity-pill severity-${rec.severity.toLowerCase()}`}>{rec.severity}</span></td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
}; 
