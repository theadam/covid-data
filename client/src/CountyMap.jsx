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

export default function CountyMap({ loading, states, counties, onDataClick }) {
  const [zoomFeature, setZoomFeature] = React.useState(null);
  const firstData = React.useMemo(() => firstArray(counties), [counties]);
  return (
    <div style={{ flex: 1 }}>
      <Map
        zoomFeature={zoomFeature}
        loading={loading}
        data={counties}
        projection={projection}
        onDataClick={onDataClick}
        formatIndex={(i) =>
          firstData[i] ? formatDate(firstData[i].date) : null
        }
        hideEmptyTip
        tipTitleFn={(_, d) => `${d.county}, ${d.state}`}
      >
        <FeatureSet
          data={counties}
          features={countyFeatures}
          getHighlight={() => true}
          dataIdKey="fipsId"
          allowEmptyDataClick
          onDataClick={() => setZoomFeature(null)}
          shouldRender={(feature) =>
            zoomFeature && feature.id.startsWith(zoomFeature.id)
          }
          calculateTip={(feature, data) => {
            if (!data) return null;
            return {
              info: covidTipInfo(data),
              title: `${data.county}, ${data.state}`,
            };
          }}
        />
        <FeatureSet
          data={states}
          features={stateFeatures}
          getHighlight={() => true}
          getStroke={() => '#888888'}
          getFill={(feature) =>
            zoomFeature && zoomFeature.id === feature.id ? 'none' : 'default'
          }
          dataIdKey="fipsId"
          calculateTip={covidTip}
          allowEmptyDataClick
          onDataClick={(_, feature) => setZoomFeature(feature)}
        />
      </Map>
    </div>
  );
}
