import React from 'react';
import '../node_modules/react-vis/dist/style.css';
import { curveCatmullRom } from 'd3-shape';

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
  Crosshair,
} from 'react-vis';
import { css } from 'emotion';

const legendClass = css`
  text-align: right;
`;

const dateRegexp = /^(\d{4})-0?(\d{1,2})-0?(\d{1,2})T/;
function formatDate(d) {
  const [, , /*year*/ month, day] = dateRegexp.exec(d);
  return `${month}/${day}`;
}

function formatNumber(d) {
  if (d < 1000) {
    return d;
  }
  return `${d / 1000}k`;
}

function formatItem(item) {
  return {
    title: item.country,
    value: `${item.confirmed} confirmed cases`,
  };
}

function formatTitle([item]) {
  return {
    title: item.formattedDate,
  };
}

function mapEachArray(obj, fn) {
  return Object.keys(obj).reduce((acc, key) => {
    return { ...acc, [key]: obj[key].map(fn) };
  }, {});
}

export default function() {
  const [data, setData] = React.useState(null);
  const [crosshairValues, setCrosshairValues] = React.useState([]);
  React.useEffect(() => {
    fetch(
      '/data/country/historical?country=United%20States,Italy,China,Spain,South%20Korea,Iran',
    )
      .then(r => r.json())
      .then(j => {
        setData(
          mapEachArray(j, (item, i) => ({
            type: 'COUNTRY',
            x: i,
            y: item.confirmed,
            formattedDate: formatDate(item.date),
            ...item,
          })),
        );
      });
  }, []);

  if (!data) return null;

  const items = Object.keys(data);

  const Plot = makeWidthFlexible(XYPlot);

  return (
    <div>
      <Plot
        height={500}
        style={{ paddingLeft: 20, paddingRight: 20, overflow: 'visible' }}
      >
        <HorizontalGridLines />
        <VerticalGridLines />
        <XAxis tickFormat={i => data[items[0]][i].formattedDate} />
        <YAxis tickFormat={formatNumber} />
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
            curve={curveCatmullRom.alpha(0.5)}
            onNearestX={(value, { index }) => {
              setCrosshairValues(items.map(i => data[i][index]));
            }}
            data={data[name]}
          />
        ))}
        <Crosshair
          values={crosshairValues}
          itemsFormat={items => items.map(formatItem)}
          titleFormat={formatTitle}
        />
      </Plot>
      <DiscreteColorLegend
        className={legendClass}
        items={items}
        orientation="horizontal"
      />
    </div>
  );
}
