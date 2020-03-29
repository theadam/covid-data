import React from 'react';
import { geoAlbersUsa as proj, geoPath } from 'd3-geo';
import { interpolateReds } from 'd3-scale-chromatic';
import * as topojson from 'topojson';
import countyData from './data/counties-10m.json';
import { interpolate } from 'd3-interpolate';

import Slider from '@material-ui/core/Slider';

import { withStyles } from '@material-ui/core/styles';

const AirbnbSlider = withStyles({
  root: {
    color: '#3a8589',
    height: 3,
    padding: '13px 0',
  },
  thumb: {
    height: 27,
    width: 27,
    backgroundColor: '#fff',
    border: '1px solid currentColor',
    marginTop: -12,
    marginLeft: -13,
    boxShadow: '#ebebeb 0px 2px 2px',
    '&:focus, &:hover, &$active': {
      boxShadow: '#ccc 0px 2px 3px 1px',
    },
    '& .bar': {
      // display: inline-block !important;
      height: 9,
      width: 1,
      backgroundColor: 'currentColor',
      marginLeft: 1,
      marginRight: 1,
    },
  },
  active: {},
  valueLabel: {
    marginLeft: 4,
    left: 'calc(-50% + 4px)',
  },
  track: {
    height: 3,
  },
  rail: {
    color: '#d8d8d8',
    opacity: 1,
    height: 3,
  },
})(Slider);

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
const countyFeatures = topojson.feature(countyData, countyData.objects.counties)
  .features;

function getIndex(data, index) {
  const keys = Object.keys(data);
  return keys.map((key) => {
    const list = data[key];
    return list[index !== undefined ? index : data[key].length - 1];
  });
}

function mapByCode(finals) {
  const result = {};
  finals.forEach((final) => {
    if (!final || !final.fipsId) return;
    result[final.fipsId] = final;
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

const dateRegexp = /^(\d{4})-0?(\d{1,2})-0?(\d{1,2})T/;
function formatDate(d) {
  const [, , /*year*/ month, day] = dateRegexp.exec(d);
  return `${month}/${day}`;
}

export default function CountyMap({ data, onDataClick }) {
  const keys = React.useMemo(() => Object.keys(data), [data]);
  const firstKey = React.useMemo(() => keys[0], [keys]);
  const firstData = React.useMemo(() => data[firstKey], [data, firstKey]);
  const [index, setIndex] = React.useState(firstData.length - 1);

  const finals = React.useMemo(() => getIndex(data), [data]);
  const dataSlice = React.useMemo(() => getIndex(data, index), [data, index]);
  const byCode = React.useMemo(() => mapByCode(dataSlice), [dataSlice]);
  const max = React.useMemo(() => getMax(finals), [finals]);

  const [tipLocation, _] = React.useState(null);
  const [width, setWidth] = React.useState(
    () => document.documentElement.clientWidth - 40,
  );
  const [height, setHeight] = React.useState(() => sizeProjection(width));
  const [path, setPath] = React.useState(() =>
    geoPath().projection(projection),
  );
  const paths = React.useMemo(() => countyFeatures.map((d) => path(d)), [path]);

  React.useEffect(() => {
    function listener() {
      const width = document.documentElement.clientWidth - 40;
      setWidth(width);
      const height = sizeProjection(width);
      setHeight(height);
      setPath(() => geoPath().projection(projection));
    }
    const debounced = debounce(listener, 400);
    window.addEventListener('resize', debounced);
    return () => window.removeEventListener('resize', debounced);
  }, []);
  const interpolator = interpolateReds;

  return (
    <div style={{ marginLeft: 20, marginRight: 20 }}>
      <svg
        height={height}
        width={width}
        style={{ border: '1px solid #AAAAAA' }}
      >
        <g height={height} width={width}>
          <g height={height} width={width}>
            {countyFeatures.map((d, i) => {
              const data = byCode[d.id];
              return (
                <path
                  key={i}
                  d={paths[i]}
                  fill={
                    data
                      ? interpolator(
                          Math.sqrt(Math.sqrt(data.confirmed)) /
                            Math.sqrt(Math.sqrt(max)),
                        )
                      : '#EEE'
                  }
                  stroke="#AAAAAA"
                  cursor={data ? 'pointer' : undefined}
                  onClick={
                    data
                      ? () => {
                          onDataClick(data);
                        }
                      : undefined
                  }
                />
              );
            })}
          </g>
        </g>
      </svg>
      {tipLocation &&
        (() => {
          const tipData = byCode[tipLocation.id];
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
                {tipData && `\n${tipData.confirmed} Confirmed Cases`}
                {tipData && `\n${tipData.deaths} Fatalities`}
                {!tipData && `\nNo Cases`}
              </div>
            </div>
          );
        })()}
      <AirbnbSlider
        valueLabelDisplay={'auto'}
        value={index}
        min={0}
        max={firstData.length - 1}
        marks
        onChange={(_, i) => setIndex(i)}
        valueLabelFormat={(i) => formatDate(firstData[i].date)}
      />
    </div>
  );
}
