import 'ol/ol.css';

import React from 'react';
import { Map } from 'ol';
import VectorSource from 'ol/source/Vector';
import VectorLayer from 'ol/layer/Vector';
import TileLayer from 'ol/layer/Tile';
import TopoJSON from 'ol/format/TopoJSON';
import countries from './data/countries-110m.json';
import { Fill, Stroke, Style } from 'ol/style';

import XYZ from 'ol/source/XYZ';

import View from 'ol/View';
const colorStart = '#FFEDA0';

export default class OpenLayersMap extends React.Component {
  componentDidMount() {
    this.map = new Map({
      target: this.el,
      layers: [
        new VectorLayer({
          source: new VectorSource({
            features: new TopoJSON({
              layers: ['countries'],
            }).readFeatures(countries, {
              dataProjection: 'EPSG:4326',
              featureProjection: 'EPSG:3857',
            }),
            overlaps: false,
          }),
          style: () => {
            return new Style({
              stroke: new Stroke({
                width: 1,
                color: 'white',
                lineDash: [3, 3],
              }),
              fill: new Fill({
                color: colorStart,
              }),
            });
          },
        }),
        new TileLayer({
          source: new XYZ({
            url:
              'https://abcd.basemaps.cartocdn.com/light_only_labels/{z}/{x}/{y}.png',
            maxZoom: 19,
          }),
        }),
      ],
      view: new View({
        center: [0, 0],
        zoom: 2,
      }),
    });
  }

  componentWillUnmount() {
    this.map.destroy();
  }

  render() {
    return (
      <div
        style={{
          background: 'rgb(202, 210, 211)',
          flex: 1,
        }}
        ref={(el) => (this.el = el)}
      />
    );
  }
}
