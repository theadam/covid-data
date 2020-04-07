import React from 'react';
import { GeoJSON } from 'react-leaflet';
import { interpolateReds } from 'd3-scale-chromatic';
import { getAllMax } from './utils';

function getVal(data, feature, index, key) {
  if (!data) return data;
  const fData = data[feature.id];
  if (!fData) return null;
  const iData = fData[index];
  if (!iData) return null;
  return iData[key];
}

function interpolate(data, max) {
  if (!data) {
    return '#eee';
  }
  return interpolateReds(
    Math.sqrt(Math.sqrt(data)) / Math.sqrt(Math.sqrt(max)),
  );
}

export default function DataLayer({
  index,
  features,
  data,
  getShow = () => true,
  getStroke,
  dataKey = 'confirmed',
  style = () => ({}),
}) {
  const max = React.useMemo(() => getAllMax(data, dataKey), [data, dataKey]);

  return (
    <GeoJSON
      data={features}
      style={(feature) => {
        const show = getShow(feature);
        const stroke = getStroke ? getStroke(feature) : show;
        const st = style(feature);
        return {
          weight: 1,
          stroke,
          color: '#AAAAAA',
          fillOpacity: 1,
          fillColor: show
            ? interpolate(getVal(data, feature, index, dataKey), max)
            : 'none',
          ...st,
        };
      }}
    />
  );
}
