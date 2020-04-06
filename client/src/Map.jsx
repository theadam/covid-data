import React from 'react';
import { css } from '@emotion/core';

import MapTip from './MapTip';
import PlaySlider from './PlaySlider';
import Loader from './Loader';
import {
  getDataIndex,
  mapBy,
  getMax,
  useProjection,
  usePlayer,
  firstArray,
  transformBounds,
} from './utils';

function getTip(tipLocation, datas, children) {
  if (!tipLocation) return null;
  const { byCode } = datas[tipLocation.featureSetIndex];
  const data = byCode[tipLocation.feature.id];
  const calculateTip = React.Children.toArray(children)[
    tipLocation.featureSetIndex
  ].props.calculateTip;
  return calculateTip(tipLocation.feature, data);
}

export default function Map({
  loading,
  data,
  projection,
  formatIndex,
  zoomFeature,
  children,
}) {
  const { ref, path, width, height } = useProjection(projection);
  const [tipLocation, setTipLocation] = React.useState(null);
  const firstData = React.useMemo(() => firstArray(data), [data]);
  const { play, playing, frame: index, setFrame: setIndex } = usePlayer(
    firstData.length,
  );

  const baseDatas = React.Children.map(children, ({ props }) => ({
    data: props.data || [],
    dataIdKey: props.dataIdKey || 'countryCode',
    dataKey: props.dataKey || 'confirmed',
  }));

  const datas = React.useMemo(
    () =>
      baseDatas.map(({ data, dataIdKey, dataKey }) => {
        const finals = getDataIndex(data);
        const dataSlice = getDataIndex(data, index);
        return {
          data,
          finals,
          dataSlice,
          byCode: mapBy(dataSlice, dataIdKey),
          max: getMax(finals, dataKey),
        };
      }),
    [baseDatas, index],
  );

  const zoomTransform = React.useMemo(() => {
    if (!zoomFeature) return { scale: 1, translate: [0, 0] };
    const bounds = path.bounds(zoomFeature);
    const dx = bounds[1][0] - bounds[0][0];
    const dy = bounds[1][1] - bounds[0][1];
    const x = (bounds[0][0] + bounds[1][0]) / 2;
    const y = (bounds[0][1] + bounds[1][1]) / 2;
    const scale = 0.7 / Math.max(dx / width, dy / height);
    const translate = [width / 2 - scale * x, height / 2 - scale * y];
    return { scale, translate };
  }, [path, zoomFeature, height, width]);

  const zoomTransformString = React.useMemo(() => {
    return `translate(${zoomTransform.translate}) scale(${zoomTransform.scale})`;
  }, [zoomTransform]);

  const tip = getTip(tipLocation, datas, children);
  const showTip = !loading && tipLocation && tip;

  return (
    <div
      css={css`
        user-select: none;
      `}
      style={{ position: 'relative' }}
      ref={ref}
    >
      <div
        className="map-container"
        style={{ border: '1px solid rgb(170, 170, 170)', display: 'flex' }}
      >
        <div style={{ position: 'relative', flex: 1 }}>
          <Loader loading={loading} />
          <svg
            height={height}
            width={width}
            style={{
              opacity: loading ? '50%' : undefined,
            }}
          >
            <g
              height={height}
              width={width}
              transform={zoomTransformString}
              style={{ transition: 'transform 0.5s' }}
            >
              <g height={height} width={width}>
                {React.Children.map(children, (child, i) => {
                  return React.cloneElement(child, {
                    path,
                    index,
                    loading,
                    onMouseOver: (feature) => {
                      if (
                        loading ||
                        !child.props.calculateTip ||
                        (tipLocation && tipLocation.id === feature.id)
                      ) {
                        return;
                      }
                      setTipLocation({
                        bounds: path.bounds(feature),
                        feature: feature,
                        featureSetIndex: i,
                      });
                    },
                    ...datas[i],
                  });
                })}
              </g>
            </g>
          </svg>
        </div>
      </div>
      {!loading && (
        <MapTip
          info={tip ? tip.info : null}
          title={tip ? tip.title : null}
          height={height}
          bounds={
            showTip ? transformBounds(tipLocation.bounds, zoomTransform) : null
          }
        />
      )}
      <PlaySlider
        playing={playing}
        play={play}
        index={index}
        length={firstData.length}
        setIndex={setIndex}
        formatLabel={formatIndex}
        hideTip={loading}
      />
    </div>
  );
}
