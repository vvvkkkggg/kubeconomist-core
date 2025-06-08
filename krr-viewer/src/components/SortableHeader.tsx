import { ArrowDown, ArrowUp, ChevronsUpDown } from 'lucide-react';
import React from 'react';
import type { SortDirection } from '../hooks/useSort';

interface SortableHeaderProps<T> {
  children: React.ReactNode;
  sortKey: keyof T;
  currentSortKey: keyof T;
  direction: SortDirection;
  onRequestSort: (key: keyof T) => void;
}

export const SortableHeader = <T extends object>({
  children,
  sortKey,
  currentSortKey,
  direction,
  onRequestSort,
}: SortableHeaderProps<T>) => {
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
