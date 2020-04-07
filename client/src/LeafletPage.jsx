import React from 'react';
import { css } from '@emotion/core';
import { Pane, Map, TileLayer } from 'react-leaflet';
import countyData from './data/counties-10m.json';
import * as topojson from 'topojson-client';
import worldData from './data/countries-50m.json';
import provinceData from './data/canadaprovtopo.json';
import PlaySlider from './PlaySlider';
import { firstArray, usePlayer, formatDate } from './utils';
import DataLayer from './DataLayer';
import Loader from './Loader';

var landMap = {
  url:
    'https://stamen-tiles-{s}.a.ssl.fastly.net/terrain-labels/{z}/{x}/{y}{r}.{ext}',
  attribution:
    'Map tiles by <a href="http://stamen.com">Stamen Design</a>, <a href="http://creativecommons.org/licenses/by/3.0">CC BY 3.0</a> &mdash; Map data &copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
  subdomains: 'abcd',
  minZoom: 0,
  maxZoom: 17,
  ext: 'png',
};

const countyFeatures = topojson.feature(
  countyData,
  countyData.objects.counties,
);

const stateFeatures = topojson.feature(countyData, countyData.objects.states);

const worldFeatures = topojson.feature(worldData, worldData.objects.countries);
const provinceFeatures = topojson.feature(
  provinceData,
  provinceData.objects.canadaprov,
);

const stateThreshold = 3;
const countyThreshold = 5;

export default function LeafletPage() {
  const [data, setData] = React.useState(null);
  React.useEffect(() => {
    Promise.all([
      fetch('/api/data/countries/historical/').then((r) => r.json()),
      fetch('/api/data/us/states/historical/').then((r) => r.json()),
      fetch('/api/data/us/counties/historical/').then((r) => r.json()),
    ]).then(([countries, states, counties]) =>
      setData({ countries, states, counties }),
    );
  }, []);
  const firstData = React.useMemo(
    () => (data ? firstArray(data.countries) : []),
    [data],
  );

  const loading = data === null;
  const { play, playing, frame: index, setFrame: setIndex } = usePlayer(
    firstData.length,
  );
  const [zoom, setZoom] = React.useState(4);
  const position = [37.0902, -95.7129];
  return (
    <div
      css={css`
        user-select: none;
        display: flex;
        flex-direction: column;
        flex: 1;
      `}
      style={{ position: 'relative' }}
    >
      <div style={{ flex: 1, display: 'flex' }}>
        <Loader loading={loading} />
        <Map
          center={position}
          zoom={zoom}
          style={{ flex: 1, opacity: loading ? 0.5 : undefined }}
          onViewportChanged={({ zoom: newZoom }) =>
            zoom !== newZoom && setZoom(newZoom)
          }
        >
          <DataLayer
            index={index}
            features={countyFeatures}
            getShow={() => zoom >= countyThreshold}
            data={data && data.counties}
          />
          <DataLayer
            index={index}
            features={stateFeatures}
            getShow={() => zoom >= stateThreshold}
            data={data && data.states}
          />
          <DataLayer
            index={index}
            features={provinceFeatures}
            getShow={() => zoom >= stateThreshold}
          />
          <DataLayer
            index={index}
            features={worldFeatures}
            data={data && data.countries}
            getShow={(feature) =>
              zoom < stateThreshold ||
              (feature.id !== '840' && feature.id !== '124')
            }
          />
          <Pane>
            <TileLayer {...landMap} noWrap zIndex={10} />
          </Pane>
        </Map>
      </div>
      <PlaySlider
        playing={playing}
        play={play}
        index={index}
        length={firstData.length}
        setIndex={setIndex}
        formatLabel={(i) =>
          firstData[i] ? formatDate(firstData[i].date) : null
        }
        hideTip={false}
      />
    </div>
  );
}
