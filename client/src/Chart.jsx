import React from 'react';
import { curveCatmullRom } from 'd3-shape';
import '../node_modules/react-vis/dist/style.css';

import {
  XYPlot,
  XAxis,
  YAxis,
  ChartLabel,
  HorizontalGridLines,
  VerticalGridLines,
  LineSeries,
  DiscreteColorLegend,
  makeWidthFlexible,
} from 'react-vis';
import { css } from 'emotion';

const legendClass = css``;

export default function() {
  const [data, setData] = React.useState(null);
  React.useEffect(() => {
    fetch('/data/country/historical?country=United%20States,Italy')
      .then(r => r.json())
      .then(setData);
  }, []);

  if (!data) return null;

  const items = Object.keys(data);

  const Plot = makeWidthFlexible(XYPlot);

  return (
    <div style={{ paddingLeft: 20, paddingRight: 20 }}>
      <Plot height={500} style={{ paddingLeft: 20, paddingRight: 20 }}>
        <HorizontalGridLines />
        <VerticalGridLines />
        <XAxis tickFormat={i => data[items[0]][i].date} tickLabelAngle={90} />
        <YAxis />
        <ChartLabel
          text="Date"
          includeMargin={false}
          xPercent={0.025}
          yPercent={1.01}
        />

        <ChartLabel
          text="Confirmed Cases"
          className="alt-y-label"
          includeMargin={false}
          xPercent={0.01}
          yPercent={0.06}
          style={{
            transform: 'rotate(-90)',
            textAnchor: 'end',
          }}
        />
        {items.map((name, i) => (
          <LineSeries
            data={data[name].map((item, i) => ({
              x: i,
              y: item.confirmed,
            }))}
          />
        ))}
      </Plot>
      <DiscreteColorLegend
        className={legendClass}
        items={items}
        orientation="horizontal"
      />
    </div>
  );
}
