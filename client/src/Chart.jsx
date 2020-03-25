import React from 'react';
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
  MarkSeries,
} from 'react-vis';
import { css } from 'emotion';
import ControlledHighlight from './ControlledHighlight';

const legendClass = css`
  text-align: right;
`;

const crosshairClass = css`
  .rv-crosshair__line {
    // z-index: -1;
    position: relative;
  }
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

function getYDomain(data, xDomain) {
  const keys = Object.keys(data);
  let maxY = null;
  let minY = null;

  keys.forEach(key => {
    data[key].slice(xDomain[0], xDomain[1] + 1).forEach(item => {
      if (maxY === null || item.y > maxY) {
        maxY = item.y;
      }
      if (minY === null || item.y < minY) {
        minY = item.y;
      }
    });
  });
  return [minY, maxY * 1.05];
}

const defaultWidth = 15;
function getInitialDomain(data) {
  const keys = Object.keys(data);
  const last = data[keys[0]].length - 1;

  if (last <= defaultWidth) {
    return { left: 0, right: last };
  }
  return { left: last - defaultWidth, right: last };
}

const Plot = makeWidthFlexible(XYPlot);

export default function() {
  const [data, setData] = React.useState(null);
  const [crosshairValues, setCrosshairValues] = React.useState([]);
  const [domain, setDomain] = React.useState(null);
  const brushing = React.useRef(false);
  React.useEffect(() => {
    fetch(
      '/data/countries/historical?country=United%20States,Italy,Spain,South%20Korea,Iran',
    )
      .then(r => r.json())
      .then(j => {
        const data = mapEachArray(j, (item, i) => ({
          x: i,
          index: i,
          y: item.confirmed,
          formattedDate: formatDate(item.date),
          ...item,
        }));
        setData(data);
        setDomain(getInitialDomain(data));
      });
  }, []);

  if (!data) return null;

  const items = Object.keys(data);

  return (
    <div style={{ paddingLeft: 20, paddingRight: 20 }}>
      <Plot
        animation
        height={500}
        style={{ overflow: 'visible' }}
        xDomain={domain && [domain.left, domain.right]}
        yDomain={domain && getYDomain(data, [domain.left, domain.right])}
      >
        <HorizontalGridLines />
        <VerticalGridLines />
        <XAxis
          tickFormat={i => {
            if (data[items[0]][i]) {
              return data[items[0]][i].formattedDate;
            }
          }}
        />
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
        <Crosshair
          className={crosshairClass}
          values={crosshairValues}
          itemsFormat={items => items.map(formatItem)}
          titleFormat={formatTitle}
        />
        {items.map((name, i) => (
          <LineSeries
            key={name}
            curve={curveCatmullRom.alpha(0.5)}
            onNearestX={
              i === 0
                ? (value, { index }) => {
                    if (!brushing.current) {
                      setCrosshairValues(items.map(i => data[i][index]));
                    }
                  }
                : null
            }
            data={data[name]}
          />
        ))}
        {crosshairValues[0]
          ? items.map(name => (
              <MarkSeries
                animation={false}
                key={name}
                stroke="white"
                data={[data[name][crosshairValues[0].index]]}
              />
            ))
          : null}
      </Plot>
      <Brusher
        setDomain={setDomain}
        brushing={brushing}
        data={data}
        items={items}
      />
      <DiscreteColorLegend
        className={legendClass}
        items={items}
        orientation="horizontal"
      />
    </div>
  );
}

const Brusher = React.memo(({ setDomain, brushing, data, items }) => {
  const [area, setArea] = React.useState(() => getInitialDomain(data));

  return (
    <Plot animation height={100} style={{ overflow: 'visible' }}>
      {items.map((name, i) => (
        <LineSeries
          key={name}
          curve={curveCatmullRom.alpha(0.5)}
          data={data[name]}
        />
      ))}
      <ControlledHighlight
        drag
        key="highlight"
        area={area}
        enableY={false}
        onBrushStart={() => {
          brushing.current = true;
        }}
        onBrushEnd={area => {
          brushing.current = false;
          setArea(area);
          setDomain(area);
        }}
        onDrag={a => {
          setDomain(a);
        }}
        onDragEnd={a => {
          setArea(a);
          setDomain(a);
        }}
      />
    </Plot>
  );
});
