import { unparse } from 'papaparse';
import React, { useMemo, useState } from 'react';
import { SortableHeader } from '../components/SortableHeader';
import { ViewHeader } from '../components/ViewHeader';
import { mockStorageRecommendations } from '../data/storage-mock-data';
import { useSort } from '../hooks/useSort';

export const StorageView: React.FC = () => {
    const [searchQuery, setSearchQuery] = useState('');

    const filteredRecs = useMemo(() => {
        const lowerCaseQuery = searchQuery.toLowerCase();
        return mockStorageRecommendations.filter(rec => {
            if (rec.type === 'Block Storage') {
                return rec.volumeName.toLowerCase().includes(lowerCaseQuery) ||
                    rec.instanceId.toLowerCase().includes(lowerCaseQuery)
            }
            return rec.bucketName.toLowerCase().includes(lowerCaseQuery)
        });
    }, [searchQuery]);

    const { items: sortedRecs, requestSort, sortKey, sortDirection } = useSort(filteredRecs, 'severity', 'descending');

    const sortProps = {
        currentSortKey: sortKey,
        direction: sortDirection,
        onRequestSort: requestSort,
    };

    const handleExport = () => {
        const dataToExport = sortedRecs.map(rec => {
            if (rec.type === 'Block Storage') {
                return {
                    Type: rec.type,
                    VolumeName: rec.volumeName,
                    InstanceId: rec.instanceId,
                    CurrentSizeGB: rec.currentSizeGB,
                    RecommendedSizeGB: rec.recommendedSizeGB,
                    Severity: rec.severity,
                }
            }
            return {
                Type: rec.type,
                BucketName: rec.bucketName,
                Region: rec.region,
                StorageClass: rec.storageClass,
                PotentialSavings: rec.potentialSavings,
                Severity: rec.severity,
            }
        });
        const csv = unparse(dataToExport);
        const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' });
        const link = document.createElement('a');
        link.href = URL.createObjectURL(blob);
        link.setAttribute('download', 'storage-recommendations.csv');
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
    };

    return (
        <div>
            <h2 className="view-header">Storage Recommendations ({sortedRecs.length})</h2>
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
                            <SortableHeader sortKey="type" {...sortProps}>Type</SortableHeader>
                            <th>Identifier</th>
                            <th>Details</th>
                            <SortableHeader sortKey="severity" {...sortProps}>Severity</SortableHeader>
                        </tr>
                    </thead>
                    <tbody>
                        {sortedRecs.map((rec) => (
                            <tr key={rec.id}>
                                <td>{rec.type}</td>
                                {rec.type === 'Block Storage' ? (
                                    <>
                                        <td>{rec.volumeName} ({rec.instanceId})</td>
                                        <td>Size: {rec.currentSizeGB}GB → <strong>{rec.recommendedSizeGB}GB</strong></td>
                                    </>
                                ) : (
                                    <>
                                        <td>{rec.bucketName} ({rec.region})</td>
                                        <td>Move from {rec.storageClass} to save <strong>${rec.potentialSavings.toFixed(2)}</strong></td>
                                    </>
                                )}
                                <td><span className={`severity-pill severity-${rec.severity.toLowerCase()}`}>{rec.severity}</span></td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
}; 
