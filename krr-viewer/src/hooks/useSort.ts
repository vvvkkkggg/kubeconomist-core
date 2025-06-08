import { useMemo, useState } from 'react';

export type SortDirection = 'ascending' | 'descending';

const SEVERITY_ORDER = { "CRITICAL": 3, "WARNING": 2, "OK": 1, "GOOD": 1, "UNKNOWN": 0 };


export const useSort = <T extends object>(
  items: T[],
  initialSortKey: keyof T,
  initialDirection: SortDirection = 'ascending'
) => {
  const [sortKey, setSortKey] = useState<keyof T>(initialSortKey);
  const [sortDirection, setSortDirection] = useState<SortDirection>(initialDirection);

  const sortedItems = useMemo(() => {
    return [...items].sort((a, b) => {
      const aValue = a[sortKey];
      const bValue = b[sortKey];

      if (aValue === null || aValue === undefined) return 1;
      if (bValue === null || bValue === undefined) return -1;
      
      if (typeof sortKey === 'string' && sortKey.toLowerCase().includes('severity')) {
          const aSeverity = SEVERITY_ORDER[aValue as keyof typeof SEVERITY_ORDER] ?? -1;
          const bSeverity = SEVERITY_ORDER[bValue as keyof typeof SEVERITY_ORDER] ?? -1;
          return sortDirection === 'ascending' ? aSeverity - bSeverity : bSeverity - aSeverity;
      }
      
      if (aValue < bValue) {
        return sortDirection === 'ascending' ? -1 : 1;
      }
      if (aValue > bValue) {
        return sortDirection === 'ascending' ? 1 : -1;
      }
      return 0;
    });
  }, [items, sortKey, sortDirection]);

  const requestSort = (key: keyof T) => {
    let direction: SortDirection = 'ascending';
    if (sortKey === key && sortDirection === 'ascending') {
      direction = 'descending';
    }
    setSortKey(key);
    setSortDirection(direction);
  };

  return { items: sortedItems, requestSort, sortKey, sortDirection };
}; 
