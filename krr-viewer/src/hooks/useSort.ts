import { useMemo, useState } from 'react';

export type SortDirection = 'ascending' | 'descending';

// Helper to resolve nested property values.
const get = (obj: any, path: string): any => {
  return path.split('.').reduce((acc, part) => acc && acc[part], obj);
};

const SEVERITY_ORDER = { "CRITICAL": 3, "WARNING": 2, "OK": 1, "GOOD": 1, "UNKNOWN": 0 };


export const useSort = <T extends object>(
  items: T[],
  initialSortKey: string,
  initialDirection: SortDirection = 'ascending'
) => {
  const [sortKey, setSortKey] = useState<string>(initialSortKey);
  const [sortDirection, setSortDirection] = useState<SortDirection>(initialDirection);

  const sortedItems = useMemo(() => {
    const sorted = [...items].sort((a, b) => {
      const aValue = get(a, sortKey);
      const bValue = get(b, sortKey);

      if (aValue === null || aValue === undefined) return 1;
      if (bValue === null || bValue === undefined) return -1;
      
      if (sortKey.toLowerCase().includes('severity')) {
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

    return sorted;
  }, [items, sortKey, sortDirection]);

  const requestSort = (key: string) => {
    let direction: SortDirection = 'ascending';
    if (sortKey === key && sortDirection === 'ascending') {
      direction = 'descending';
    }
    setSortKey(key);
    setSortDirection(direction);
  };

  return { items: sortedItems, requestSort, sortKey, sortDirection };
}; 
