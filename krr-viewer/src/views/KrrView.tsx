import { Download, Search } from 'lucide-react';
import { unparse } from 'papaparse';
import React, { useMemo, useState } from 'react';
import { RecommendationsTable } from '../components/RecommendationsTable';
import { mockReports } from '../data/mock-data';
import { useSort } from '../hooks/useSort';
import type { KrrReport, Scan } from '../types';

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

const KrrHeader: React.FC<{
    onSearch: (query: string) => void;
    onHideEmptyToggle: (hidden: boolean) => void;
    scans: Scan[];
    reports: KrrReport[];
    selectedReport: KrrReport;
    onReportChange: (report: KrrReport) => void;
}> = ({ onSearch, onHideEmptyToggle, scans, reports, selectedReport, onReportChange }) => {

    const handleExport = () => {
        const dataToExport = scans.map(s => ({
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

    return (
        <div className="krr-header">
            <div className="controls-left">
                <div className="search-container">
                    <Search size={18} />
                    <input type="text" placeholder="Type to search" onChange={(e) => onSearch(e.target.value)} />
                </div>
                <label className="checkbox-container">
                    <input type="checkbox" onChange={(e) => onHideEmptyToggle(e.target.checked)} />
                    Hide empty recommendations
                </label>
            </div>
            <div className="controls-right">
                <select className="scan-date-select" value={selectedReport.date} onChange={(e) => onReportChange(reports.find(r => r.date === e.target.value)!)}>
                    {reports.map(r => <option key={r.date} value={r.date}>Scan date: {formatDate(r.date)}</option>)}
                </select>
                <button className="export-button" onClick={handleExport}>
                    <Download size={16} />
                    Export CSV
                </button>
            </div>
        </div>
    );
};


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


    return (
        <div>
            <h2 className="view-header">Recommendations ({sortedScans.length})</h2>
            <KrrHeader
                onSearch={setSearchQuery}
                onHideEmptyToggle={setHideEmpty}
                scans={sortedScans}
                reports={mockReports}
                selectedReport={selectedReport}
                onReportChange={setSelectedReport}
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
