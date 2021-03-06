import React from 'react';
import { GeoJSON } from 'react-leaflet';
// import { interpolateReds } from 'd3-scale-chromatic';
import { interpolate as interp } from 'd3-interpolate';
import { color } from 'd3-color';

const colorStart = '#FFEDA0';
const colorEnd = '#800026';
const interpolateReds = interp(colorStart, colorEnd);

function interpolate(data, max) {
  if (data === null || data === undefined) {
    return color(colorStart).brighter(0.2).formatHex();
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
    style = () => ({}),
    onSelect,
    calculateValue,
  }) => {
    const propsRef = React.useRef({ index, data, style });
    propsRef.current = { index, data, style };

    return (
      <GeoJSON
        data={featureCollection}
        onEachFeature={(feature, layer) => {
          layer.on({
            mouseover: () => {
              layer.bringToFront();
              layer.setStyle({
                weight: 3,
                dashArray: '',
              });
              onHighlight({
                feature: feature,
                data: propsRef.current?.data?.[feature.key],
                displayName: feature.displayName,
              });
            },
            mouseout: () => {
              layer.setStyle({
                weight: 1,
                dashArray: 3,
                ...propsRef.current?.style(
                  feature,
                  propsRef.current?.data?.[feature.key]?.[
                    propsRef.current?.index
                  ],
                ),
              });
              onHighlight(null);
            },
            click: () => {
              onSelect(propsRef.current?.data?.[feature.key]);
            },
          });
        }}
        style={(feature) => {
          const show = getShow(feature);
          const stroke = getStroke ? getStroke(feature) : show;
          const item = data?.[feature.key];
          const dates = item?.dates;
          const date = dates?.[index];
          const value = calculateValue(item, index);
          const st = style(feature, item, date);
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
