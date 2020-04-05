import React from 'react';
import MapPath from './MapPath';

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
  path,
  index,
  loading,
  byCode,
  max,
  onMouseOver,
}) {
  const svgPaths = React.useMemo(() => features.map((d) => path(d)), [
    features,
    path,
  ]);

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
            onMouseOver={() => onMouseOver(feature)}
          />
        );
      })}
    </>
  );
}
