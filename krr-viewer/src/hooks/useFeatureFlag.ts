import { useMemo } from 'react';

export const useFeatureFlag = (flagName: string): boolean => {
  const isEnabled = useMemo(
    () => import.meta.env[flagName] === 'true',
    [flagName]
  );
  return isEnabled;
}; 
