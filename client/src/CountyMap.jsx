import React from 'react';
import { geoAlbersUsa as proj } from 'd3-geo';
import * as topojson from 'topojson';
import countyData from './data/counties-10m.json';

import { covidTip, covidTipInfo, formatDate, firstArray } from './utils';

import Map from './Map';
import FeatureSet from './FeatureSet';

const projection = proj();

const countyFeatures = topojson.feature(countyData, countyData.objects.counties)
  .features;

const stateFeatures = topojson.feature(countyData, countyData.objects.states)
  .features;

export default function CountyMap({ loading, data, onDataClick }) {
  const firstData = React.useMemo(() => firstArray(data), [data]);
  return (
    <div style={{ flex: 1 }}>
      <Map
        loading={loading}
        data={data}
        projection={projection}
        onDataClick={onDataClick}
        formatIndex={(i) =>
          firstData[i] ? formatDate(firstData[i].date) : null
        }
        hideEmptyTip
        tipTitleFn={(_, d) => `${d.county}, ${d.state}`}
      >
        <FeatureSet
          data={data}
          features={countyFeatures}
          getHighlight={() => true}
          dataIdKey="fipsId"
          calculateTip={(feature, data) => {
            if (!data) return null;
            return {
              info: covidTipInfo(data),
              title: `${data.county}, ${data.state}`,
            };
          }}
        />
        <FeatureSet
          features={stateFeatures}
          getHighlight={() => false}
          getStroke={() => '#888888'}
          getFill={() => 'none'}
          dataIdKey="fipsId"
          calculateTip={covidTip}
        />
      </Map>
    </div>
  );
}
