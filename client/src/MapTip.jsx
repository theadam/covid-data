import React from 'react';

export default function MapTip({ bounds, height, title, info, emptyText }) {
  const [[x1, y1], [x2, y2]] = bounds;
  const yOffset = y1 + (y2 - y1) / 2 > height / 2 ? -(height / 6) : height / 6;

  return (
    <div
      style={{
        position: 'absolute',
        left: 0,
        top: 0,
        transform: `translate(${x1 + (x2 - x1) / 2}px,${
          (y2 > height / 2 ? y1 : y2) + yOffset
        }px`,
        transition: 'transform 0.3s',
        pointerEvents: 'none',
        whiteSpace: 'pre',
      }}
      className="chart-tip"
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
        {info && info.map((val, k) => <span key={k}>{`\n${val}`}</span>)}
        {!info && `\n${emptyText}`}
      </div>
    </div>
  );
}
