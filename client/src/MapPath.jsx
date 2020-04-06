import React from 'react';
import { interpolateReds } from 'd3-scale-chromatic';
import { css } from '@emotion/core';

export default function MapPath({
  path,
  data,
  max,
  onClick,
  topoData,
  onMouseOver,
  stroke,
  fill,
  highlightOpacity,
  highlight,
}) {
  if (!topoData) return null;
  return (
    <path
      css={css`
        &.highlight:hover {
          opacity: ${highlightOpacity};
        }
      `}
      className={highlight !== false ? 'highlight' : null}
      key={topoData.id}
      d={path}
      fill={
        fill && fill !== 'default'
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
