import React from 'react';
import { geoAlbersUsa as proj } from 'd3-geo';
import * as topojson from 'topojson-client';
import countyData from './data/counties-10m.json';
import IconButton from '@material-ui/core/IconButton';
import ZoomIn from '@material-ui/icons/ZoomIn';
import ZoomOut from '@material-ui/icons/ZoomOut';
import Cancel from '@material-ui/icons/Cancel';

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
  const [inZoomMode, setZoomMode] = React.useState(false);
  const firstData = React.useMemo(() => firstArray(counties), [counties]);
  return (
    <div style={{ flex: 1, position: 'relative' }}>
      <div style={{ position: 'absolute', top: 0, right: 0, zIndex: 1000 }}>
        {inZoomMode ? (
          <span style={{ fontSize: 8, marginRight: -10 }}>
            Click on a state...
          </span>
        ) : undefined}
        <IconButton
          onClick={() => {
            if (zoomFeature) {
              return setZoomFeature(null);
            } else if (inZoomMode) {
              return setZoomMode(false);
            } else {
              return setZoomMode(true);
            }
          }}
        >
          {zoomFeature ? <ZoomOut /> : inZoomMode ? <Cancel /> : <ZoomIn />}
        </IconButton>
      </div>
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
          shouldRender={(feature, _, tipLocation) =>
            (zoomFeature && feature.id.startsWith(zoomFeature.id)) ||
            (tipLocation && feature.id.startsWith(tipLocation.feature.id))
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
          getCursor={() => (inZoomMode ? 'zoom-in' : 'default')}
          highlightOpacity={0.2}
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
          onDataClick={
            inZoomMode
              ? (_, feature) => {
                  setZoomFeature(feature);
                  setZoomMode(false);
                }
              : null
          }
        />
      </Map>
    </div>
  );
}
