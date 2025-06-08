import React from 'react';

interface InstanceGroupLinkProps {
  folderId: string;
  nodeGroupId: string;
}

export const InstanceGroupLink: React.FC<InstanceGroupLinkProps> = ({ folderId, nodeGroupId }) => {
  if (!folderId || !nodeGroupId) {
    return <span>{nodeGroupId}</span>;
  }

  const href = `https://console.yandex.cloud/folders/${folderId}/compute/instance-group/${nodeGroupId}/overview`;

  return (
    <a href={href} target="_blank" rel="noopener noreferrer">
      {nodeGroupId}
    </a>
  );
}; 
