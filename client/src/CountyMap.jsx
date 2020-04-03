import React from 'react';
import { geoAlbersUsa as proj } from 'd3-geo';
import * as topojson from 'topojson';
import countyData from './data/counties-10m.json';

import { formatDate, firstArray } from './utils';

import Map from './Map';

const projection = proj();

const countyFeatures = topojson.feature(countyData, countyData.objects.counties)
  .features;

export default function CountyMap({ loading, data, onDataClick }) {
  const firstData = React.useMemo(() => firstArray(data), [data]);

  return (
    <div style={{ flex: 1 }}>
      <Map
        loading={loading}
        data={data}
        projection={projection}
        features={countyFeatures}
        onDataClick={onDataClick}
        formatIndex={(i) => formatDate(firstData[i].date)}
        dataIdKey="fipsId"
        hideEmptyTip
        tipTitleFn={(_, d) => `${d.county}, ${d.state}`}
      />
    </div>
  );
}
