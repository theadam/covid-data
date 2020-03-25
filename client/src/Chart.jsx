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

  keys.forEach((key) => {
    data[key].slice(xDomain[0], xDomain[1] + 1).forEach((item) => {
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

const defaultWidth = 30;
function getInitialDomain(data) {
  const keys = Object.keys(data);
  if (keys.length === 0) return undefined;
  const last = data[keys[0]].length - 1;

  if (last <= defaultWidth) {
    return { left: 0, right: last };
  }
  return { left: last - defaultWidth, right: last };
}

const Plot = makeWidthFlexible(XYPlot);

function pick(map, items) {
  if (items.length === 0) {
    return makeWorldData(map);
  }
  return items.reduce((acc, k) => {
    return { ...acc, [k]: map[k] };
  }, {});
}

function makeWorldData(data) {
  const keys = Object.keys(data);

  const result = [];
  const base = () => ({
    confirmed: 0,
    deaths: 0,
    date: '',
    counry: '',
  });
  keys.forEach((key) => {
    data[key].forEach((d, i) => {
      if (!result[i]) {
        result[i] = base();
        result[i].date = d.date;
        result[i].country = 'World';
      }
      result[i].confirmed += d.confirmed;
      result[i].deaths += d.deaths;
    });
  });
  return { World: result };
}

const overrides = {};
const DAY = 1000 * 60 * 60 * 24;

function moveDate(date, override) {
  if (override === 0) {
    return date;
  }
  const [, year, month, day] = dateRegexp.exec(date);

  const parsed = new Date(
    Date.UTC(Number(year), Number(month) - 1, Number(day)),
  );
  parsed.setTime(parsed.getTime() + override * DAY);
  return parsed.toISOString();
}

export default function ({ data: baseData, chartedCountries, onLegendClick }) {
  const [crosshairValues, setCrosshairValues] = React.useState([]);
  const brushing = React.useRef(false);
  const data = React.useMemo(
    () =>
      mapEachArray(pick(baseData, chartedCountries), (item, i) => {
        const override = overrides[item.country] || 0;
        const date = moveDate(item.date, override);

        return {
          x: i + override,
          index: i + override,
          y: item.confirmed,
          formattedDate: formatDate(date),
          date,
          ...item,
        };
      }),
    [baseData, chartedCountries],
  );
  const [domain, setDomain] = React.useState(() => getInitialDomain(data));

  const items = Object.keys(data);

  return (
    <div style={{ paddingLeft: 20, paddingRight: 20, flex: 1 }}>
      <div>
        <Plot
          animation
          height={300}
          xDomain={domain && [domain.left, domain.right]}
          yDomain={domain && getYDomain(data, [domain.left, domain.right])}
        >
          <HorizontalGridLines />
          <VerticalGridLines />
          <XAxis
            tickFormat={(i) => {
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
            itemsFormat={(items) => items.map(formatItem)}
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
                        setCrosshairValues(items.map((i) => data[i][index]));
                      }
                    }
                  : null
              }
              data={data[name]}
            />
          ))}
          {crosshairValues[0]
            ? items.map((name) => (
                <MarkSeries
                  animation={false}
                  key={name}
                  stroke="white"
                  data={[data[name][crosshairValues[0].index]]}
                />
              ))
            : null}
        </Plot>
      </div>
      <div>
        <Brusher
          setDomain={setDomain}
          brushing={brushing}
          data={data}
          items={items}
        />
      </div>
      <DiscreteColorLegend
        onItemClick={chartedCountries.length > 0 ? onLegendClick : undefined}
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
        onBrushEnd={(area) => {
          brushing.current = false;
          setArea(area);
          setDomain(area);
        }}
        onDrag={(a) => {
          setDomain(a);
        }}
        onDragEnd={(a) => {
          setArea(a);
          setDomain(a);
        }}
      />
    </Plot>
  );
});
