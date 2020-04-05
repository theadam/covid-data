import React from 'react';
import { css } from '@emotion/core';

import MapTip from './MapTip';
import PlaySlider from './PlaySlider';
import Loader from './Loader';
import { useProjection, usePlayer, firstArray, transformBounds } from './utils';

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

  const zoomTransform = React.useMemo(() => {
    if (!zoomFeature) return { scale: 1, translate: [0, 0] };
    const bounds = path.bounds(zoomFeature);
    const dx = bounds[1][0] - bounds[0][0];
    const dy = bounds[1][1] - bounds[0][1];
    const x = (bounds[0][0] + bounds[1][0]) / 2;
    const y = (bounds[0][1] + bounds[1][1]) / 2;
    const scale = 0.9 / Math.max(dx / width, dy / height);
    const translate = [width / 2 - scale * x, height / 2 - scale * y];
    return { scale, translate };
  }, [path, zoomFeature, height, width]);

  const zoomTransformString = React.useMemo(() => {
    return `translate(${zoomTransform.translate}) scale(${zoomTransform.scale})`;
  }, [zoomTransform]);

  const tip = tipLocation ? tipLocation.tip : null;
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
            <g height={height} width={width} transform={zoomTransformString}>
              <g height={height} width={width}>
                {React.Children.map(children, (child) => {
                  return React.cloneElement(child, {
                    path,
                    tipLocation,
                    setTipLocation,
                    index,
                    loading,
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
