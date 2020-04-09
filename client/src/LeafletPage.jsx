import React from 'react';
import { css } from '@emotion/core';
import { Pane, Map, TileLayer } from 'react-leaflet';
import countyData from './data/counties-10m.json';
import * as topojson from 'topojson-client';
import worldData from './data/countries-50m.json';
import provinceData from './data/canadaprovtopo.json';
import australiaData from './data/au-states.json';
import chinaData from './data/china-provinces.json';
import PlaySlider from './PlaySlider';
import { getAllMax, firstArray, usePlayer, formatDate } from './utils';
import DataLayer from './DataLayer';
import Loader from './Loader';
import fipsData from './fipsData.json';

const USA = '840';
const Canada = '124';
const China = '156';
const Australia = '036';

const countriesWithRegions = [USA, Canada, China, Australia];

var landMap = {
  ext: 'png',
  url: 'https://{s}.basemaps.cartocdn.com/light_only_labels/{z}/{x}/{y}{r}.png',
  attribution:
    '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors &copy; <a href="https://carto.com/attributions">CARTO</a>',
  subdomains: 'abcd',
  maxZoom: 19,
};

const stateThreshold = 3;
const countyThreshold = 5;

const baseFeatures = {
  world: topojson.feature(worldData, worldData.objects.countries),
  usStates: topojson.feature(countyData, countyData.objects.states),
  usCounties: topojson.feature(countyData, countyData.objects.counties),
  china: topojson.feature(chinaData, chinaData.objects.CHN_adm1),
  australia: topojson.feature(australiaData, australiaData.objects.states),
  canada: topojson.feature(provinceData, provinceData.objects.canadaprov),
};

function enrichFeatures(
  featureCollection,
  featureKey = (feature) => feature.id,
  displayName = (feature) => feature.properties.name,
) {
  return {
    ...featureCollection,
    features: featureCollection.features.map((feature) => {
      const key = featureKey(feature);
      return {
        ...feature,
        key,
        displayName: displayName(feature),
      };
    }),
  };
}

function createFeatures() {
  return {
    world: enrichFeatures(baseFeatures.world),
    usStates: enrichFeatures(baseFeatures.usStates),
    usCounties: enrichFeatures(
      baseFeatures.usCounties,
      undefined,
      (feature) => {
        if (!fipsData[feature.id]) {
          console.log(feature.id);
        }
        return fipsData[feature.id].displayName;
      },
    ),
    china: enrichFeatures(
      baseFeatures.china,
      (feature) => `${China}-${feature.properties.NAME_1}`,
      (feature) => feature.properties.NAME_1,
    ),
    australia: enrichFeatures(
      baseFeatures.australia,
      (feature) => `${Australia}-${feature.properties.STATE_NAME}`,
      (feature) => feature.properties.STATE_NAME,
    ),
    canada: enrichFeatures(
      baseFeatures.canada,
      (feature) => `${Canada}-${feature.properties.name}`,
    ),
  };
}

const features = createFeatures();

function pluckFeatureData(featureCollection, data) {
  return featureCollection.features.reduce((acc, { key }) => {
    if (data[key]) {
      acc[key] = data[key];
    }
    return acc;
  }, {});
}

function splitData(data) {
  return {
    world: pluckFeatureData(features.world, data.countries),
    usStates: pluckFeatureData(features.usStates, data.states),
    usCounties: pluckFeatureData(features.usCounties, data.counties),
    china: pluckFeatureData(features.china, data.provinces),
    australia: pluckFeatureData(features.australia, data.provinces),
    canada: pluckFeatureData(features.canada, data.provinces),
  };
}

const dataKey = 'confirmed';

export default function LeafletPage() {
  const [data, setData] = React.useState({});
  React.useEffect(() => {
    fetch('/api/data/all/historical/')
      .then((r) => r.json())
      .then((data) => {
        setData(splitData(data));
      });
  }, []);

  const [highlight, setHighlight] = React.useState(null);
  const firstData = React.useMemo(() => (data ? firstArray(data.world) : []), [
    data,
  ]);
  const [zoom, setZoom] = React.useState(4);
  const showCounties = React.useMemo(() => zoom >= countyThreshold, [zoom]);
  const showProvinces = React.useMemo(() => zoom >= stateThreshold, [zoom]);
  const getShowCounties = React.useCallback(() => showCounties, [showCounties]);
  const getShowProvinces = React.useCallback(() => showProvinces, [
    showProvinces,
  ]);

  const provinceMax = React.useMemo(
    () => getAllMax({ ...data.provinces, ...data.usStates }, dataKey),
    [data],
  );
  const countiesMax = React.useMemo(() => getAllMax(data.usCounties, dataKey), [
    data,
  ]);
  const worldMax = React.useMemo(() => getAllMax(data.world, dataKey), [data]);

  const loading = !data.world;
  const { play, playing, frame: index, setFrame: setIndex } = usePlayer(
    firstData.length,
  );
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
          worldCopyJump
          center={position}
          zoom={zoom}
          style={{
            flex: 1,
            opacity: loading ? 0.5 : undefined,
            background: 'rgb(202, 210, 211)',
            position: 'relative',
          }}
          onViewportChanged={({ zoom: newZoom }) =>
            zoom !== newZoom && setZoom(newZoom)
          }
        >
          {!loading && (
            <div
              className="highlight-info"
              css={css`
                z-index: 1000;
                position: absolute;
                top: 10px;
                right: 10px;
                padding: 6px 8px;
                font: 14px/16px Arial, Helvetica, sans-serif;
                background: white;
                background: rgba(255, 255, 255, 0.8);
                box-shadow: 0 0 15px rgba(0, 0, 0, 0.2);
                border-radius: 5px;
                h4 {
                  margin: 0 0 5px;
                  color: #777;
                }
              `}
            >
              {highlight === null ? (
                <h4>Hover over an area to see information</h4>
              ) : (
                <span>
                  <b>{highlight.displayName}</b>
                  <br />
                  {highlight?.dataArray?.[index] ? (
                    <span>
                      <span>
                        {highlight?.dataArray?.[
                          index
                        ]?.confirmed.toLocaleString()}{' '}
                        Confirmed Cases
                        <br />
                      </span>
                      <span>
                        {highlight?.dataArray?.[index]?.deaths.toLocaleString()}{' '}
                        Deaths
                        <br />
                      </span>
                    </span>
                  ) : (
                    <span>No Cases</span>
                  )}
                </span>
              )}
            </div>
          )}
          <DataLayer
            index={index}
            name="usCounties"
            featureCollection={features.usCounties}
            getShow={getShowCounties}
            data={data.usCounties}
            onHighlight={setHighlight}
            max={countiesMax}
          />
          <DataLayer
            index={index}
            name="usStates"
            featureCollection={features.usStates}
            data={data.usStates}
            getShow={getShowProvinces}
            onHighlight={setHighlight}
            max={provinceMax}
            style={React.useCallback(
              () =>
                showCounties
                  ? {
                      weight: 2,
                      fillColor: 'none',
                    }
                  : {},
              [showCounties],
            )}
          />
          <DataLayer
            index={index}
            name="canada"
            featureCollection={features.canada}
            data={data.canada}
            getShow={getShowProvinces}
            onHighlight={setHighlight}
            max={provinceMax}
          />
          <DataLayer
            index={index}
            name="china"
            featureCollection={features.china}
            data={data.china}
            getShow={getShowProvinces}
            onHighlight={setHighlight}
            max={provinceMax}
          />
          <DataLayer
            index={index}
            name="australia"
            featureCollection={features.australia}
            data={data.australia}
            getShow={getShowProvinces}
            onHighlight={setHighlight}
            max={provinceMax}
          />
          <DataLayer
            index={index}
            name="world"
            featureCollection={features.world}
            data={data.world}
            onHighlight={setHighlight}
            max={worldMax}
            style={React.useCallback(
              (feature) =>
                showProvinces
                  ? {
                      weight: 2,
                      ...(countriesWithRegions.includes(feature.id)
                        ? { fillColor: 'none' }
                        : {}),
                    }
                  : {},
              [showProvinces],
            )}
          />
          <Pane>
            {React.useMemo(
              () => (
                <TileLayer {...landMap} noWrap />
              ),
              [],
            )}
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
