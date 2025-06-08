import { useEffect, useState } from 'react';
import { krrReport as mockKrrReport } from '../data/krr-mock-data';
import type { KrrReport } from '../types';
import { useFeatureFlag } from './useFeatureFlag';

export const useKrrData = () => {
  const useBackend = useFeatureFlag('VITE_USE_BACKEND_DATA');

  const [data, setData] = useState<KrrReport | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      setError(null);
      try {
        if (useBackend) {
          const response = await fetch('/api/krr'); // Or whatever the real endpoint is
          if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
          }
          const result = await response.json();
          setData(result);
        } else {
          setData(mockKrrReport);
        }
      } catch (e) {
        if (e instanceof Error) {
            setError(e);
        } else {
            setError(new Error('An unknown error occurred'));
        }
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [useBackend]);

  return { data, loading, error };
}; 
