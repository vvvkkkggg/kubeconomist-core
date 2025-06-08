import { unparse } from 'papaparse';
import React, { useMemo, useState } from 'react';
import { SortableHeader } from '../components/SortableHeader';
import { ViewHeader } from '../components/ViewHeader';
import { mockVpcRecommendations } from '../data/vpc-mock-data';
import { useSort } from '../hooks/useSort';
import type { VpcRecommendation } from '../types';

const FolderLink: React.FC<{ folderId: string }> = ({ folderId }) => (
    <a href={`https://console.yandex.cloud/folders/${folderId}`} target="_blank" rel="noopener noreferrer">
      {folderId}
    </a>
);

const VpcTable: React.FC<{ items: VpcRecommendation[] }> = ({ items }) => {
    const { items: sortedRecs, requestSort, sortKey: currentSortKey, sortDirection } = useSort(items, 'ipAddress', 'ascending');
    const sortProps = { currentSortKey, direction: sortDirection, onRequestSort: requestSort };

    return (
        <div className="table-container">
            <table>
                <thead>
                    <tr>
                        <SortableHeader sortKey="cloudId" {...sortProps}>Cloud ID</SortableHeader>
                        <SortableHeader sortKey="folderId" {...sortProps}>Folder ID</SortableHeader>
                        <SortableHeader sortKey="ipAddress" {...sortProps}>IP Address</SortableHeader>
                        <SortableHeader sortKey="isUsed" {...sortProps}>Is Used</SortableHeader>
                        <SortableHeader sortKey="isReserved" {...sortProps}>Is Reserved</SortableHeader>
                    </tr>
                </thead>
                <tbody>
                    {sortedRecs.map((rec) => (
                        <tr key={rec.id}>
                            <td>{rec.cloudId}</td>
                            <td><FolderLink folderId={rec.folderId} /></td>
                            <td>{rec.ipAddress}</td>
                            <td>{rec.isUsed ? 'Yes' : 'No'}</td>
                            <td>{rec.isReserved ? 'Yes' : 'No'}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
}

export const VpcView: React.FC = () => {
    const [searchQuery, setSearchQuery] = useState('');

    const filteredRecs = useMemo(() => {
        const lowerCaseQuery = searchQuery.toLowerCase();
        // The user only cares about reserved IPs that are not being used.
        return mockVpcRecommendations.filter(rec =>
            rec.isReserved && !rec.isUsed && (
                rec.cloudId.toLowerCase().includes(lowerCaseQuery) ||
                rec.folderId.toLowerCase().includes(lowerCaseQuery) ||
                rec.ipAddress.toLowerCase().includes(lowerCaseQuery)
            )
        );
    }, [searchQuery]);

    const handleExport = () => {
        const dataToExport = filteredRecs.map(rec => ({
            'Cloud ID': rec.cloudId,
            'Folder ID': rec.folderId,
            'IP Address': rec.ipAddress,
            'Is Used': rec.isUsed ? 'Yes' : 'No',
            'Is Reserved': rec.isReserved ? 'Yes' : 'No',
        }));

        const csv = unparse(dataToExport);
        const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' });
        const link = document.createElement('a');
        link.href = URL.createObjectURL(blob);
        link.setAttribute('download', `vpc-unused-ips-report.csv`);
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
    };

    return (
        <div>
            <h2 className="view-header">Unused Reserved IP Addresses ({filteredRecs.length})</h2>
            <ViewHeader
                onSearch={setSearchQuery}
                onHideEmptyToggle={() => {}}
                showHideEmpty={false}
                onExport={handleExport}
            />
            <VpcTable items={filteredRecs} />
        </div>
    );
}; 
