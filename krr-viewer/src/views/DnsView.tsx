import { unparse } from 'papaparse';
import React, { useMemo, useState } from 'react';
import { SortableHeader } from '../components/SortableHeader';
import { ViewHeader } from '../components/ViewHeader';
import { mockDnsRecommendations } from '../data/dns-mock-data';
import { useSort } from '../hooks/useSort';
import type { DnsRecommendation } from '../types';

const FolderLink: React.FC<{ folderId: string }> = ({ folderId }) => (
    <a href={`https://console.yandex.cloud/folders/${folderId}`} target="_blank" rel="noopener noreferrer">
        {folderId}
    </a>
);

export const DnsView: React.FC = () => {
    const [searchQuery, setSearchQuery] = useState('');

    const unusedRecs = useMemo(() => {
        return mockDnsRecommendations.filter((rec: DnsRecommendation) => !rec.isUsed);
    }, []);

    const filteredRecs = useMemo(() => {
        const lowerCaseQuery = searchQuery.toLowerCase();
        return unusedRecs.filter((rec: DnsRecommendation) =>
            rec.zoneId.toLowerCase().includes(lowerCaseQuery) ||
            rec.cloudId.toLowerCase().includes(lowerCaseQuery) ||
            rec.folderId.toLowerCase().includes(lowerCaseQuery)
        );
    }, [searchQuery, unusedRecs]);

    const { items: sortedRecs, requestSort, sortKey, sortDirection } = useSort<DnsRecommendation>(filteredRecs, 'zoneId', 'ascending');

    const sortProps = {
        currentSortKey: sortKey,
        direction: sortDirection,
        onRequestSort: requestSort,
    };

    const handleExport = () => {
        const dataToExport = sortedRecs.map((rec: DnsRecommendation) => ({
            'Zone ID': rec.zoneId,
            'Cloud ID': rec.cloudId,
            'Folder ID': rec.folderId,
        }));
        const csv = unparse(dataToExport);
        const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' });
        const link = document.createElement('a');
        link.href = URL.createObjectURL(blob);
        link.setAttribute('download', 'dns-recommendations.csv');
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
    };

    return (
        <div>
            <h2 className="view-header">Unused DNS Zones ({sortedRecs.length})</h2>
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
                            <SortableHeader<DnsRecommendation> sortKey="zoneId" {...sortProps}>Zone ID</SortableHeader>
                            <SortableHeader<DnsRecommendation> sortKey="cloudId" {...sortProps}>Cloud ID</SortableHeader>
                            <SortableHeader<DnsRecommendation> sortKey="folderId" {...sortProps}>Folder ID</SortableHeader>
                        </tr>
                    </thead>
                    <tbody>
                        {sortedRecs.map((rec: DnsRecommendation) => (
                            <tr key={rec.id}>
                                <td>{rec.zoneId}</td>
                                <td>{rec.cloudId}</td>
                                <td><FolderLink folderId={rec.folderId} /></td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
}; 
