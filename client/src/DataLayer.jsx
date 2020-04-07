import React from 'react';
import { GeoJSON } from 'react-leaflet';
import { interpolateReds } from 'd3-scale-chromatic';
import { getAllMax } from './utils';

function getData(data, featureKey, index) {
  if (!data) return data;
  const fData = data[featureKey];
  if (!fData) return null;
  const iData = fData[index];
  return iData;
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
  featureKey = (feature) => feature.id,
}) {
  const max = React.useMemo(() => getAllMax(data, dataKey), [data, dataKey]);

  return (
    <GeoJSON
      data={features}
      style={(feature) => {
        const show = getShow(feature);
        const stroke = getStroke ? getStroke(feature) : show;
        const item = getData(data, featureKey(feature), index);
        const st = style(feature, item);
        return {
          weight: 1,
          stroke,
          color: '#AAAAAA',
          fillOpacity: 1,
          fillColor: show
            ? interpolate(item ? item[dataKey] : null, max)
            : 'none',
          ...st,
        };
      }}
    />
  );
}
