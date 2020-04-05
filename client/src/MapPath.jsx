import React from 'react';
import { interpolateReds } from 'd3-scale-chromatic';

export default function MapPath({
  path,
  data,
  max,
  onClick,
  topoData,
  onMouseOver,
  stroke,
  fill,
  highlight,
}) {
  if (!topoData) return null;
  return (
    <path
      className={highlight !== false ? 'highlight' : null}
      key={topoData.id}
      d={path}
      fill={
        fill
          ? fill
          : data
          ? interpolateReds(
              Math.sqrt(Math.sqrt(data)) / Math.sqrt(Math.sqrt(max)),
            )
          : '#EEE'
      }
      strokeWidth={0.3}
      stroke={stroke || '#AAAAAA'}
      cursor={onClick && data ? 'pointer' : undefined}
      onClick={onClick}
      onMouseOver={onMouseOver}
    />
  );
}
