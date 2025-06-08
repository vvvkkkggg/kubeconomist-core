import { unparse } from 'papaparse';
import React, { useMemo, useState } from 'react';
import { RecommendationsTable } from '../components/RecommendationsTable';
import { ViewHeader } from '../components/ViewHeader';
import * as format from '../formatters';
import { useKrrData } from '../hooks/useKrrData';
import { useSort } from '../hooks/useSort';
import type { Scan } from '../types';

// function formatDate(iso: string) {
//     return new Date(iso).toLocaleString(undefined, {
//         year: '2-digit',
//         month: 'short',
//         day: '2-digit',
//         hour: '2-digit',
//         minute: '2-digit',
//         second: '2-digit',
//     });
// }

export const KrrView: React.FC = () => {
    const { data: krrReport, loading, error } = useKrrData();
    const [searchQuery, setSearchQuery] = useState('');
    const [hideEmpty, setHideEmpty] = useState(false);
    const [showSavingsOnly, setShowSavingsOnly] = useState(true);

    const scans = useMemo(() => {
        if (!krrReport) {
            return [];
        }

        let scans = krrReport.scans;
        if (hideEmpty) {
            scans = scans.filter(s => s.severity !== 'UNKNOWN' && s.recommended.requests.cpu.value !== '?');
        }
        if (showSavingsOnly) {
            scans = scans.filter(s => {
                const cpuReq = s.recommended.requests.cpu.value;
                const memReq = s.recommended.requests.memory.value;
                if (typeof cpuReq === 'number' && typeof s.object.allocations.requests.cpu === 'number' && cpuReq < s.object.allocations.requests.cpu) {
                    return true;
                }
                if (typeof memReq === 'number' && typeof s.object.allocations.requests.memory === 'number' && memReq < s.object.allocations.requests.memory) {
                    return true;
                }
                return false;
            })
        }
        if (searchQuery) {
            scans = scans.filter(s => s.object.container.toLowerCase().includes(searchQuery.toLowerCase()))
        }
        return scans;
    }, [krrReport, searchQuery, hideEmpty, showSavingsOnly]);

    const { items: sortedScans, requestSort, sortKey, sortDirection } = useSort<Scan>(scans, 'severity', 'descending');

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
            'Cost Savings (₽/month)': format.calculateCostSavings(s),
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

    if (loading) {
        return <div>Loading...</div>;
    }

    if (error) {
        return <div>Error: {error.message}</div>;
    }

    return (
        <div>
            <h2 className="view-header">Recommendations ({sortedScans.length})</h2>
            <ViewHeader
                onSearch={setSearchQuery}
                onHideEmptyToggle={setHideEmpty}
                onExport={handleExport}
                onShowSavingsOnlyToggle={setShowSavingsOnly}
                showSavingsOnly={true}
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
