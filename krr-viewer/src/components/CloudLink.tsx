import React from 'react';

interface CloudLinkProps {
  cloudId: string;
}

export const CloudLink: React.FC<CloudLinkProps> = ({ cloudId }) => {
  if (!cloudId) {
    return null;
  }

  const href = `https://console.yandex.cloud/cloud/${cloudId}`;

  return (
    <a href={href} target="_blank" rel="noopener noreferrer">
      {cloudId}
    </a>
  );
}; 
