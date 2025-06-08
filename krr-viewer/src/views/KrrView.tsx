import { unparse } from 'papaparse';
import React, { useMemo, useState } from 'react';
import { RecommendationsTable } from '../components/RecommendationsTable';
import { ViewHeader } from '../components/ViewHeader';
import { mockReports } from '../data/mock-data';
import { useSort } from '../hooks/useSort';
import type { KrrReport } from '../types';

function formatDate(iso: string) {
    return new Date(iso).toLocaleString(undefined, {
        year: '2-digit',
        month: 'short',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
    });
}

export const KrrView: React.FC = () => {
    const [selectedReport, setSelectedReport] = useState<KrrReport>(mockReports[0]);
    const [searchQuery, setSearchQuery] = useState('');
    const [hideEmpty, setHideEmpty] = useState(false);

    const scans = useMemo(() => {
        let scans = selectedReport.scans;
        if (hideEmpty) {
            scans = scans.filter(s => s.severity !== 'UNKNOWN' && s.recommended.requests.cpu.value !== '?');
        }
        if (searchQuery) {
            scans = scans.filter(s => s.object.container.toLowerCase().includes(searchQuery.toLowerCase()))
        }
        return scans;
    }, [selectedReport, searchQuery, hideEmpty]);

    const { items: sortedScans, requestSort, sortKey, sortDirection } = useSort(scans, 'severity', 'descending');

    const handleExport = () => {
        const dataToExport = sortedScans.map(s => ({
            Name: s.object.name,
            Namespace: s.object.namespace,
            Kind: s.object.kind,
            Container: s.object.container,
            Severity: s.severity,
            'Mem (Req)': `${s.object.allocations.requests.memory} -> ${s.recommended.requests.memory.value}`,
            'Mem (Lim)': `${s.object.allocations.limits.memory} -> ${s.recommended.limits.memory.value}`,
            'CPU (Req)': `${s.object.allocations.requests.cpu} -> ${s.recommended.requests.cpu.value}`,
            'CPU (Lim)': `${s.object.allocations.limits.cpu} -> ${s.recommended.limits.cpu.value}`,
        }));
        const csv = unparse(dataToExport);
        const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' });
        const link = document.createElement('a');
        link.href = URL.createObjectURL(blob);
        link.setAttribute('download', 'krr-report.csv');
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
    }

    const scanDateSelect = (
        <select className="scan-date-select" value={selectedReport.date} onChange={(e) => setSelectedReport(mockReports.find(r => r.date === e.target.value)!)}>
            {mockReports.map(r => <option key={r.date} value={r.date}>Scan date: {formatDate(r.date)}</option>)}
        </select>
    );

    return (
        <div>
            <h2 className="view-header">Recommendations ({sortedScans.length})</h2>
            <ViewHeader
                onSearch={setSearchQuery}
                onHideEmptyToggle={setHideEmpty}
                onExport={handleExport}
                scanDateSelect={scanDateSelect}
            />
            <RecommendationsTable
                scans={sortedScans}
                requestSort={requestSort}
                sortKey={sortKey}
                sortDirection={sortDirection}
            />
        </div>
    );
}; 
