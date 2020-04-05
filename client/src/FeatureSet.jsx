import React from 'react';
import { v4 as uuidv4 } from 'uuid';
import MapPath from './MapPath';
import { usePrevious, getDataIndex, mapBy, getMax } from './utils';

export default function FeatureSet({
  data = {},
  features,
  shouldRender,
  calculateTip,
  onDataClick,
  allowEmptyDataClick = false,
  getFill,
  getStroke,
  getHighlight,
  dataKey = 'confirmed',
  dataIdKey = 'countryCode',

  // Injected by map
  tipLocation,
  setTipLocation,
  path,
  index,
  loading,
}) {
  const featureSetUUID = React.useMemo(() => uuidv4(), []);
  const finals = React.useMemo(() => getDataIndex(data), [data]);
  const dataSlice = React.useMemo(() => getDataIndex(data, index), [
    data,
    index,
  ]);
  const byCode = React.useMemo(() => mapBy(dataSlice, dataIdKey), [
    dataSlice,
    dataIdKey,
  ]);
  const max = React.useMemo(() => getMax(finals, dataKey), [finals, dataKey]);

  const svgPaths = React.useMemo(() => features.map((d) => path(d)), [
    features,
    path,
  ]);

  const previousByCode = usePrevious(byCode);
  React.useEffect(() => {
    if (
      tipLocation &&
      tipLocation.setId === featureSetUUID &&
      previousByCode !== byCode &&
      previousByCode !== null
    ) {
      setTipLocation({
        ...tipLocation,
        data: byCode[tipLocation.id],
      });
    }
  }, [byCode, setTipLocation, tipLocation, previousByCode, featureSetUUID]);

  return (
    <>
      {features.map((feature, i) => {
        const data = byCode[feature.id];
        if (shouldRender && !shouldRender(feature, data)) {
          return null;
        }
        return (
          <MapPath
            key={i}
            highlight={getHighlight ? getHighlight(feature, data) : false}
            stroke={getStroke ? getStroke(feature, data) : null}
            fill={getFill ? getFill(feature, data) : null}
            path={svgPaths[i]}
            data={data ? data.confirmed : null}
            topoData={feature}
            max={max}
            onClick={
              onDataClick && (data || allowEmptyDataClick)
                ? () => onDataClick(data)
                : undefined
            }
            onMouseOver={(e) => {
              if (
                loading ||
                !calculateTip ||
                (tipLocation && tipLocation.id === feature.id)
              ) {
                return;
              }
              setTipLocation({
                bounds: path.bounds(feature),
                id: feature.id,
                feature,
                data: byCode[feature.id],
                tip: calculateTip(feature, data),
                setId: featureSetUUID,
              });
            }}
          />
        );
      })}
    </>
  );
}
