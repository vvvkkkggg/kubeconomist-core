import { ArrowDown, ArrowUp, ChevronsUpDown } from 'lucide-react';
import React from 'react';
import type { SortDirection } from '../hooks/useSort';

interface SortableHeaderProps {
  children: React.ReactNode;
  sortKey: string;
  currentSortKey: string;
  direction: SortDirection;
  onRequestSort: (key: string) => void;
}

export const SortableHeader: React.FC<SortableHeaderProps> = ({
  children,
  sortKey,
  currentSortKey,
  direction,
  onRequestSort,
}) => {
  const isSorting = currentSortKey === sortKey;
  const SortIcon = isSorting
    ? direction === 'ascending'
      ? ArrowUp
      : ArrowDown
    : ChevronsUpDown;

  return (
    <th onClick={() => onRequestSort(sortKey)} >
        <div className="sortable-header">
            <span>{children}</span>
            <SortIcon size={16} className={`sort-icon ${isSorting ? 'active' : ''}`} />
        </div>
    </th>
  );
}; 
