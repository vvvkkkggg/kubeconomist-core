import React from 'react';
import * as format from '../formatters';
import type { SortDirection } from '../hooks/useSort';
import type { Scan } from '../types';
import { ChangeCell } from './ChangeCell';
import { SortableHeader } from './SortableHeader';

const SeverityPill: React.FC<{ severity: string }> = ({ severity }) => {
  const lowerSeverity = severity.toLowerCase();
  return <span className={`severity-pill severity-${lowerSeverity}`}>{severity}</span>;
};

interface RecommendationsTableProps {
  scans: Scan[];
  requestSort: (key: string) => void;
  sortKey: string;
  sortDirection: SortDirection;
}

export const RecommendationsTable: React.FC<RecommendationsTableProps> = ({
  scans,
  requestSort,
  sortKey,
  sortDirection,
}) => {
  const sortProps = {
    currentSortKey: sortKey,
    direction: sortDirection,
    onRequestSort: requestSort,
  };

  return (
    <div className="table-container">
      <table>
        <thead>
          <tr>
            <SortableHeader sortKey="object.name" {...sortProps}>Name</SortableHeader>
            <SortableHeader sortKey="object.namespace" {...sortProps}>Namespace</SortableHeader>
            <SortableHeader sortKey="object.kind" {...sortProps}>Kind</SortableHeader>
            <SortableHeader sortKey="object.container" {...sortProps}>Container</SortableHeader>
            <SortableHeader sortKey="severity" {...sortProps}>Severity</SortableHeader>
            <th>Mem (Req)</th>
            <th>Mem (Lim)</th>
            <th>CPU (Req)</th>
            <th>CPU (Lim)</th>
          </tr>
        </thead>
        <tbody>
          {scans.map((scan: Scan, index: number) => (
            <tr key={`${scan.object.name}-${scan.object.container}-${index}`}>
              <td>{scan.object.name}</td>
              <td>{scan.object.namespace}</td>
              <td>{scan.object.kind}</td>
              <td>{scan.object.container}</td>
              <td>
                <SeverityPill severity={scan.severity} />
              </td>
              <td>
                <ChangeCell
                  current={scan.object.allocations.requests.memory}
                  recommended={scan.recommended.requests.memory.value}
                  formatter={format.formatMemory}
                />
              </td>
              <td>
                <ChangeCell
                  current={scan.object.allocations.limits.memory}
                  recommended={scan.recommended.limits.memory.value}
                  formatter={format.formatMemory}
                />
              </td>
              <td>
                <ChangeCell
                  current={scan.object.allocations.requests.cpu}
                  recommended={scan.recommended.requests.cpu.value}
                  formatter={format.formatCPU}
                />
              </td>
              <td>
                <ChangeCell
                  current={scan.object.allocations.limits.cpu}
                  recommended={scan.recommended.limits.cpu.value}
                  formatter={format.formatCPU}
                />
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}; 
