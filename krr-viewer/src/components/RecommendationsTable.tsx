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

const PodLink: React.FC<{ scan: Scan, children: React.ReactNode }> = ({ scan, children }) => {
    const folderId = "b1g8enqi42dal6t489ou"; // This should ideally come from a config or context
    const url = `https://console.yandex.cloud/folders/${folderId}/managed-kubernetes/cluster/${scan.object.clusterId}/pod?id=${scan.object.namespace}:${scan.object.name}`;
    return <a href={url} target="_blank" rel="noopener noreferrer" className="btn">{children}</a>;
}

interface RecommendationsTableProps<T extends Scan> {
  scans: T[];
  requestSort: (key: keyof T) => void;
  sortKey: keyof T;
  sortDirection: SortDirection;
}

export const RecommendationsTable = <T extends Scan>({
  scans,
  requestSort,
  sortKey,
  sortDirection,
}: RecommendationsTableProps<T>) => {
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
            <SortableHeader<T> sortKey={'object.name' as keyof T} {...sortProps}>Name</SortableHeader>
            <SortableHeader<T> sortKey={'object.namespace' as keyof T} {...sortProps}>Namespace</SortableHeader>
            <SortableHeader<T> sortKey={'object.kind' as keyof T} {...sortProps}>Kind</SortableHeader>
            <SortableHeader<T> sortKey={'object.container' as keyof T} {...sortProps}>Container</SortableHeader>
            <SortableHeader<T> sortKey={'severity' as keyof T} {...sortProps}>Severity</SortableHeader>
            <th>Mem (Req)</th>
            <th>Mem (Lim)</th>
            <th>CPU (Req)</th>
            <th>CPU (Lim)</th>
            <th>Экономия (₽/мес)</th>
            <th>extra</th>
          </tr>
        </thead>
        <tbody>
          {scans.map((scan: T, index: number) => (
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
              <td>
                {format.formatRubles(format.calculateCostSavings(scan))}
              </td>
              <td>
                <PodLink scan={scan}>Сделать дешевле</PodLink>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}; 
