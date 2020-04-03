import React from 'react';
import { geoNaturalEarth1 as proj } from 'd3-geo';
import * as topojson from 'topojson';
import worldData from './data/countries-110m.json';

import Map from './Map';
import { formatDate, firstArray } from './utils';

const projection = proj();

const worldFeatures = topojson.feature(
  topojson.simplify(topojson.presimplify(worldData)),
  worldData.objects.countries,
).features;

export default function WorldMap({ data, onDataClick, loading }) {
  const firstData = React.useMemo(() => firstArray(data), [data]);

  return (
    <div className="world-map" style={{ flex: 1 }}>
      <Map
        data={data}
        loading={loading}
        projection={projection}
        features={worldFeatures}
        onDataClick={onDataClick}
        formatIndex={(i) =>
          firstData[i] ? formatDate(firstData[i].date) : null
        }
      />
    </div>
  );
}
