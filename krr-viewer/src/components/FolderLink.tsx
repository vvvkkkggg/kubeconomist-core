import React from 'react';

interface FolderLinkProps {
  folderId: string;
}

export const FolderLink: React.FC<FolderLinkProps> = ({ folderId }) => {
  if (!folderId) {
    return null;
  }

  const href = `https://console.yandex.cloud/folders/${folderId}`;

  return (
    <a href={href} target="_blank" rel="noopener noreferrer">
      {folderId}
    </a>
  );
}; 
