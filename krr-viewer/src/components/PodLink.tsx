import { Scan } from '@/types';
import { Button } from '@mantine/core';
import Link from 'next/link';

export const PodLink = ({ scan, namespace, pod }: { scan: Scan; namespace: string; pod: string }) => {
  if (!scan.clusterId) {
    return null;
  }

  const href = `https://console.cloud.yandex.ru/managed-kubernetes/cluster/${scan.clusterId}/workloads?filter=name%3D${pod}&active-tab=pods&ns=${namespace}`;

  return (
    <Button component={Link} href={href} target="_blank">
      Изменить
    </Button>
  );
}; 
