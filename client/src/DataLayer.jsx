import React from 'react';
import { GeoJSON } from 'react-leaflet';
// import { interpolateReds } from 'd3-scale-chromatic';
import { interpolate as interp } from 'd3-interpolate';
import { color } from 'd3-color';
import { getAllMax } from './utils';

const colorStart = '#FFEDA0';
const colorEnd = '#800026';
const interpolateReds = interp(colorStart, colorEnd);

function interpolate(data, max) {
  if (!data) {
    return colorStart;
  }
  return color(
    interpolateReds(Math.sqrt(Math.sqrt(data)) / Math.sqrt(Math.sqrt(max))),
  )
    .brighter(0.2)
    .formatHex();
}

export default React.memo(
  ({
    index,
    data,
    featureCollection,
    getShow = () => true,
    getStroke,
    onHighlight,
    max,
    dataKey = 'confirmed',
    style = () => ({}),
  }) => {
    const dataRef = React.useRef(null);
    dataRef.current = data;

    return (
      <GeoJSON
        data={featureCollection}
        onEachFeature={(feature, layer) => {
          layer.on({
            mouseover: () => {
              onHighlight({
                dataArray: dataRef.current?.[feature.key],
                displayName: feature.displayName,
              });
            },
            mouseout: () => onHighlight(null),
          });
        }}
        style={(feature) => {
          const show = getShow(feature);
          const stroke = getStroke ? getStroke(feature) : show;
          const array = data?.[feature.key];
          const item = array?.[index];
          const value = item?.[dataKey];
          const st = style(feature, item);
          return {
            weight: 1,
            stroke,
            color: 'white',
            dashArray: 3,
            fillOpacity: 1,
            fillColor: show ? interpolate(value, max) : 'none',
            ...st,
          };
        }}
      />
    );
  },
);
