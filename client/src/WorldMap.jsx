import React from 'react';
import { geoNaturalEarth1 as proj, geoPath } from 'd3-geo';
import { interpolateReds } from 'd3-scale-chromatic';
import * as topojson from 'topojson';
import worldData from './data/countries-110m.json';

let cancel;
const debounce = (fn, time) => (...args) => {
  if (cancel) clearTimeout(cancel);
  cancel = setTimeout(() => fn(...args), time);
};

const projection = proj();
function sizeProjection(width) {
  const outline = { type: 'Sphere' };

  const [[x0, y0], [x1, y1]] = geoPath(
    projection.fitWidth(width, outline),
  ).bounds(outline);

  const dy = Math.ceil(y1 - y0),
    l = Math.min(Math.ceil(x1 - x0), dy);
  projection.scale((projection.scale() * (l - 1)) / l).precision(0.2);
  return dy;
}
sizeProjection();
const worldFeatures = topojson.feature(worldData, worldData.objects.countries)
  .features;

function getFinals(data) {
  const keys = Object.keys(data);
  return keys.map((key) => {
    const list = data[key];
    return list[list.length - 1];
  });
}

function mapByCode(finals) {
  const result = {};
  finals.forEach((final) => {
    if (!final.countryCode) return;
    result[final.countryCode] = final;
  });
  return result;
}

function getMax(finals, item = 'confirmed') {
  let max = 0;

  finals.forEach((datum) => {
    const val = datum[item];
    if (val > max) {
      max = val;
    }
  });
  return max;
}

export default function WorldMap({ data, onDataClick }) {
  const finals = getFinals(data);
  const byCode = mapByCode(finals);
  const max = getMax(finals);

  const [tipLocation, setTipLocation] = React.useState(null);
  const [width, setWidth] = React.useState(
    Math.min(document.documentElement.clientWidth, 700),
  );
  const [height, setHeight] = React.useState(() => sizeProjection(width));
  const [path, setPath] = React.useState(() =>
    geoPath().projection(projection),
  );
  React.useEffect(() => {
    function listener() {
      const width = Math.min(document.documentElement.clientWidth, 700);
      setWidth(width);
      const height = sizeProjection(width);
      setHeight(height);
      setPath(() => geoPath().projection(projection));
    }
    const debounced = debounce(listener, 400);
    window.addEventListener('resize', debounced);
    return () => window.removeEventListener('resize', debounced);
  }, []);

  return (
    <div>
      <svg height={height} width={width}>
        <g height={height} width={width}>
          <g height={height} width={width}>
            {worldFeatures.map((d, i) => {
              const data = byCode[d.id];
              return (
                <path
                  key={i}
                  d={path(d)}
                  fill={data ? interpolateReds(data.confirmed / max) : '#EEE'}
                  stroke="#AAAAAA"
                  cursor={data ? 'pointer' : undefined}
                  onClick={
                    data
                      ? () => {
                          onDataClick(data);
                        }
                      : undefined
                  }
                  onMouseOver={(e) => {
                    if (tipLocation && tipLocation.id === d.id) return;
                    setTipLocation({
                      bounds: path.bounds(d),
                      id: d.id,
                      name: d.properties.name,
                      data,
                    });
                  }}
                />
              );
            })}
          </g>
        </g>
      </svg>
      {tipLocation &&
        (() => {
          const [[x1, y1], [x2, y2]] = tipLocation.bounds;
          const yOffset =
            y1 + (y2 - y1) / 2 > height / 2 ? -(height / 6) : height / 6;

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
                <span style={{ fontWeight: 'bold' }}>{tipLocation.name}</span>
                {tipLocation.data &&
                  `\n${tipLocation.data.confirmed} Confirmed Cases`}
                {tipLocation.data && `\n${tipLocation.data.deaths} Fatalities`}
                {!tipLocation.data && `\nNo Cases`}
              </div>
            </div>
          );
        })()}
    </div>
  );
}
