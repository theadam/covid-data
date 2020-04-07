import React from 'react';
import { geoNaturalEarth1 as proj } from 'd3-geo';
import * as topojson from 'topojson-client';
import worldData from './data/countries-50m.json';

import Map from './Map';
import FeatureSet from './FeatureSet';
import { covidTipIncludingNoCases, formatDate, firstArray } from './utils';

const projection = proj();
const worldFeatures = topojson.feature(worldData, worldData.objects.countries)
  .features;

export default function WorldMap({ data, onDataClick, loading }) {
  const firstData = React.useMemo(() => firstArray(data), [data]);

  return (
    <div className="world-map" style={{ flex: 1 }}>
      <Map
        loading={loading}
        data={data}
        projection={projection}
        formatIndex={(i) =>
          firstData[i] ? formatDate(firstData[i].date) : null
        }
      >
        <FeatureSet
          data={data}
          features={worldFeatures}
          calculateTip={covidTipIncludingNoCases}
          onDataClick={onDataClick}
          getHighlight={() => true}
        />
      </Map>
    </div>
  );
}
