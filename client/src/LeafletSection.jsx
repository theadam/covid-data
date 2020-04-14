import React from 'react';
import { css } from '@emotion/core';
import { Pane, Map, TileLayer } from 'react-leaflet';
import { getAllMax, makePolyline } from './utils';
import DataLayer from './DataLayer';
import Loader from './Loader';
import features, { data, countriesWithRegions } from './features';

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

const dataKey = 'confirmed';
const defaultPosition = [0, 0];
const defaultZoom = 2;

export default function LeafletPage({ centeredItem, onSelect, index }) {
  const mapRef = React.useRef();
  const [highlight, setHighlight] = React.useState(null);
  React.useEffect(() => {
    if (centeredItem) {
      mapRef.current.leafletElement.fitBounds(
        makePolyline(centeredItem.geometry).getBounds(),
      );
    } else {
      mapRef.current.leafletElement.setView(defaultPosition, defaultZoom);
    }
  }, [centeredItem]);
  const [zoom, setZoom] = React.useState(defaultZoom);
  const showCounties = React.useMemo(() => zoom >= countyThreshold, [zoom]);
  const showProvinces = React.useMemo(() => zoom >= stateThreshold, [zoom]);
  const getShowCounties = React.useCallback(() => showCounties, [showCounties]);
  const getShowProvinces = React.useCallback(() => showProvinces, [
    showProvinces,
  ]);

  const provinceMax = React.useMemo(
    () => getAllMax({ ...data.provinces, ...data.usStates }, dataKey),
    [],
  );
  const countiesMax = React.useMemo(
    () => getAllMax(data.usCounties, dataKey),
    [],
  );
  const worldMax = React.useMemo(() => getAllMax(data.world, dataKey), []);

  const loading = !data.world;
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
      <div
        style={{
          flex: 1,
          flexBasis: 600,
          maxHeight: '65vh',
          display: 'flex',
        }}
      >
        <Loader loading={loading} />
        <Map
          ref={mapRef}
          worldCopyJump
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
                pointer-events: none;
                z-index: 1000;
                position: absolute;
                top: 10px;
                right: 10px;
                padding: 6px 8px;
                font: 14px/16px Arial, Helvetica, sans-serif;
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
                  {highlight?.data?.dates?.[index] ? (
                    <span>
                      <span>
                        {highlight?.data?.dates?.[
                          index
                        ]?.confirmed.toLocaleString()}{' '}
                        Confirmed Cases
                        <br />
                      </span>
                      <span>
                        {highlight?.data?.dates?.[
                          index
                        ]?.deaths.toLocaleString()}{' '}
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
            onSelect={onSelect}
            index={index}
            name="usCounties"
            featureCollection={features.usCounties}
            getShow={getShowCounties}
            data={data.usCounties}
            onHighlight={setHighlight}
            max={countiesMax}
          />
          <DataLayer
            onSelect={onSelect}
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
            onSelect={onSelect}
            index={index}
            name="canada"
            featureCollection={features.canada}
            data={data.canada}
            getShow={getShowProvinces}
            onHighlight={setHighlight}
            max={provinceMax}
          />
          <DataLayer
            onSelect={onSelect}
            index={index}
            name="china"
            featureCollection={features.china}
            data={data.china}
            getShow={getShowProvinces}
            onHighlight={setHighlight}
            max={provinceMax}
          />
          <DataLayer
            onSelect={onSelect}
            index={index}
            name="australia"
            featureCollection={features.australia}
            data={data.australia}
            getShow={getShowProvinces}
            onHighlight={setHighlight}
            max={provinceMax}
          />
          <DataLayer
            onSelect={onSelect}
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
    </div>
  );
}
