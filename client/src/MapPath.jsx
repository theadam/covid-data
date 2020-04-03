import React from 'react';
import { interpolateReds } from 'd3-scale-chromatic';

export default function MapPath({
  path,
  data,
  max,
  onDataClick,
  topoData,
  onMouseOver,
}) {
  if (!topoData) return null;
  return (
    <path
      key={topoData.id}
      d={path}
      fill={
        data
          ? interpolateReds(
              Math.sqrt(Math.sqrt(data)) / Math.sqrt(Math.sqrt(max)),
            )
          : '#EEE'
      }
      strokeWidth={0.3}
      stroke="#AAAAAA"
      cursor={onDataClick && data ? 'pointer' : undefined}
      onClick={
        onDataClick && data
          ? () => {
              onDataClick(data);
            }
          : undefined
      }
      onMouseOver={onMouseOver}
    />
  );
}
