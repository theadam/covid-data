import React from 'react';

import MapPath from './MapPath';
import MapTip from './MapTip';
import PlaySlider from './PlaySlider';
import Loader from './Loader';
import {
  useProjection,
  getDataIndex,
  mapBy,
  getMax,
  usePlayer,
  firstArray,
} from './utils';

export default function Map({
  loading,
  data,
  projection,
  features,
  onDataClick,
  formatIndex,
  hideEmptyTip = false,
  tipInfo = (d) =>
    d && [`${d.confirmed} Confirmed Cases`, `${d.deaths} Fatalities`],
  tipTitleFn = (tipLocation) => tipLocation.name,
  dataItem = 'confirmed',
  dataIdKey = 'countryCode',
}) {
  const { ref, path, height } = useProjection(projection);
  const [tipLocation, setTipLocation] = React.useState(null);
  const firstData = React.useMemo(() => firstArray(data), [data]);
  const { play, playing, frame: index, setFrame: setIndex } = usePlayer(
    firstData.length,
  );

  const finals = React.useMemo(() => getDataIndex(data), [data]);
  const dataSlice = React.useMemo(() => getDataIndex(data, index), [
    data,
    index,
  ]);
  const byCode = React.useMemo(() => mapBy(dataSlice, dataIdKey), [
    dataSlice,
    dataIdKey,
  ]);
  const max = React.useMemo(() => getMax(finals, dataItem), [finals, dataItem]);
  const svgPaths = React.useMemo(() => features.map((d) => path(d)), [
    features,
    path,
  ]);

  const tipData = tipLocation ? byCode[tipLocation.id] : null;

  return (
    <div style={{ position: 'relative' }} ref={ref}>
      <div style={{ position: 'relative' }}>
        <Loader loading={loading} />
        <svg
          height={height}
          width="100%"
          style={{
            border: '1px solid #AAAAAA',
            opacity: loading ? '50%' : undefined,
          }}
        >
          <g height={height} width="100%">
            <g height={height} width="100%">
              {features.map((d, i) => {
                const data = byCode[d.id];
                return (
                  <MapPath
                    key={i}
                    path={svgPaths[i]}
                    data={data ? data.confirmed : null}
                    topoData={d}
                    max={max}
                    onDataClick={onDataClick}
                    onMouseOver={(e) => {
                      if (tipLocation && tipLocation.id === d.id) return;
                      setTipLocation({
                        bounds: path.bounds(d),
                        id: d.id,
                        name: d.properties.name,
                      });
                    }}
                  />
                );
              })}
            </g>
          </g>
        </svg>
      </div>
      {tipLocation && (!hideEmptyTip || tipData) && (
        <MapTip
          hideEmptyTip
          emptyText="No Cases"
          info={tipInfo(tipData)}
          title={tipTitleFn(tipLocation, tipData)}
          height={height}
          bounds={tipLocation.bounds}
        />
      )}
      <PlaySlider
        playing={playing}
        play={play}
        index={index}
        length={firstData.length}
        setIndex={setIndex}
        valueLabelFormat={formatIndex}
        hideTip={loading}
      />
    </div>
  );
}
