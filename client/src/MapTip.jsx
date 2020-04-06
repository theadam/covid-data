import React from 'react';

function calculateTransform(bounds, height) {
  if (bounds === null) return null;
  const [[x1, y1], [x2, y2]] = bounds;
  const midY = y1 + (y2 - y1) / 2;
  const yOffset = midY > height / 2 ? -(height / 6) : height / 6;
  return `translate(${x1 + (x2 - x1) / 2}px,${
    (midY > height / 2 ? y1 : y2) + yOffset
  }px`;
}

export default function MapTip({ bounds, height, title, info }) {
  const transform = calculateTransform(bounds, height);

  return (
    <div
      style={{
        display: !info ? 'none' : undefined,
        position: 'absolute',
        left: 0,
        top: 0,
        transform,
        transition: 'transform 0.3s',
        pointerEvents: 'none',
      }}
      className="map-tip"
    >
      <div
        style={{
          transform: `translate(-50%, -50%)`,
          borderRadius: 4,
          background: '#3a3a48',
          color: '#fff',
          fontSize: 12,
          padding: '7px 10px',
          boxShadow: '0 2px 4px rgba(0,0,0,0.5)',
        }}
      >
        <span style={{ fontWeight: 'bold' }}>{title}</span>
        {info}
      </div>
    </div>
  );
}
