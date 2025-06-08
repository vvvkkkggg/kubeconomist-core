import React from 'react';

interface InstanceLinkProps {
  folderId: string;
  instanceId: string;
}

export const InstanceLink: React.FC<InstanceLinkProps> = ({ folderId, instanceId }) => {
  if (!folderId || !instanceId) {
    return <span>{instanceId}</span>;
  }

  const href = `https://console.yandex.cloud/folders/${folderId}/compute/instance/${instanceId}/overview`;

  return (
    <a href={href} target="_blank" rel="noopener noreferrer">
      {instanceId}
    </a>
  );
}; 
