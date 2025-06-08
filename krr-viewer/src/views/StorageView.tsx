import React from 'react';
import { SortableHeader } from '../components/SortableHeader';
import { mockStorageRecommendations } from '../data/storage-mock-data';
import { useSort } from '../hooks/useSort';

export const StorageView: React.FC = () => {
    const { items: sortedRecs, requestSort, sortKey, sortDirection } = useSort(mockStorageRecommendations, 'severity', 'descending');

    const sortProps = {
        currentSortKey: sortKey,
        direction: sortDirection,
        onRequestSort: requestSort,
    };

    return (
        <div>
            <h2 className="view-header">Storage Recommendations ({sortedRecs.length})</h2>
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
