import { ArrowDown, ArrowUp } from 'lucide-react';
import React from 'react';
import { formatChange } from '../formatters';

interface ChangeCellProps {
  current: number | null | string;
  recommended: number | null | string;
  formatter: (value: number | null | string) => string;
}

export const ChangeCell: React.FC<ChangeCellProps> = ({ current, recommended, formatter }) => {
  const isIncrease = (typeof recommended === 'number' && typeof current === 'number' && recommended > current) ||
                     (typeof recommended === 'number' && current === null);

  const isDecrease = (typeof recommended === 'number' && typeof current === 'number' && recommended < current) ||
                     (current !== null && recommended === null);

  const Arrow = isIncrease ? ArrowUp : isDecrease ? ArrowDown : null;
  const color = isIncrease ? '#e57373' : isDecrease ? '#81c784' : 'inherit';

  const formattedValue = formatChange(current, recommended, formatter);
  const [oldVal, newVal] = formattedValue.split('→');

  return (
    <div style={{ display: 'flex', alignItems: 'center' }}>
      {Arrow && <Arrow size={16} style={{ color, marginRight: '4px' }} />}
      <span>
        {oldVal}→ <strong>{newVal}</strong>
      </span>
    </div>
  );
}; 
