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
import { css as globalCss } from 'emotion';
import ControlledHighlight from './ControlledHighlight';
import Loader from './Loader';
import { worldItem, dateRange, allData } from './features';
import { formatDate, values } from './utils';

function calculateIncrease(current, old) {
  if (old === null || old === undefined) return undefined;
  return Math.round(((current - old) / current) * 1000) / 10;
}

const mapCalculations = {
  confirmed: (item) => item.confirmed,
  deaths: (item) => item.deaths,
  increaseConfirmed: (item, i, items) =>
    calculateIncrease(item.confirmed, items?.[i - 1]?.confirmed),
  increaseDeaths: (item, i, items) =>
    calculateIncrease(item.confirmed, items?.[i - 1]?.confirmed),
};

export const typeText = {
  confirmed: 'Confirmed Cases',
  deaths: 'Deaths',
  increaseConfirmed: '% Change in Cases',
  increaseDeaths: '% Change in Deaths',
};

const formatValue = {
  confirmed: (n) => `${n.toLocaleString()} Confirmed Cases`,
  deaths: (n) => `${n.toLocaleString()} Deaths`,
  increaseConfirmed: (n) => `${n.toLocaleString()}% Growth of Cases`,
  increaseDeaths: (n) => `${n.toLocaleString()}% Growth of Deaths`,
};

const legendClass = globalCss`
  text-align: right;
`;

const noSelect = globalCss`
  user-select: none;
`;

const crosshairClass = globalCss`
  .rv-crosshair__line {
    // z-index: -1;
    position: relative;
  }
`;

function formatNumber(d) {
  if (d < 1000) {
    return d;
  }
  if (d < 1000000) {
    return `${d / 1000}k`;
  }
  return `${d / 1000000}m`;
}

function formatItem(formatValue, item) {
  return {
    title: item.displayName,
    value: formatValue(item.y),
  };
}

function formatTitle([item]) {
  return {
    title: formatDate(item.date),
  };
}

function mapObject(obj, fn) {
  return Object.keys(obj).reduce((acc, key) => {
    return { ...acc, [key]: fn(obj[key]) };
  }, {});
}

function compact(a) {
  return a.filter((i) => !!i);
}

function removeLastEach(obj) {
  return Object.keys(obj).reduce((acc, key) => {
    return {
      ...acc,
      [key]: { ...obj[key], dates: obj[key].dates.slice(0, -1) },
    };
  }, {});
}

function getYDomain(data, xDomain) {
  const keys = Object.keys(data);
  let maxY = null;
  let minY = null;

  keys.forEach((key) => {
    data[key].dates.slice(xDomain[0], xDomain[1] + 1).forEach((item) => {
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
  const last = data[keys[0]].dates.length - 1;

  if (last <= defaultWidth) {
    return { left: 0, right: last };
  }
  return { left: last - defaultWidth, right: last };
}

const Plot = makeWidthFlexible(XYPlot);

function pick(map, items) {
  if (items.length === 0) {
    return { [worldItem.key]: worldItem };
  }
  return items.reduce((acc, k) => {
    return { ...acc, [k]: map[k] };
  }, {});
}

export default function ({
  selected,
  loading,
  onLegendClick,
  type = 'increaseConfirmed',
}) {
  const [crosshairValues, setCrosshairValues] = React.useState([]);
  const brushing = React.useRef(false);
  const formattedDates = dateRange.map(formatDate);
  const formatter = React.useMemo(() => formatValue[type], [type]);
  const text = React.useMemo(() => typeText[type], [type]);
  const data = React.useMemo(() => {
    const calc = mapCalculations[type];
    return removeLastEach(
      mapObject(
        pick(
          allData,
          selected.map((item) => item.key),
        ),
        (item) => {
          return {
            ...item,
            dates: compact(
              item.dates.map((date, i, dates) => {
                const y = calc(date, i, dates);
                return {
                  ...item,
                  ...date,
                  x: i,
                  index: i,
                  y: y || 0,
                };
              }),
            ),
          };
        },
      ),
    );
  }, [selected, type]);
  const [domain, rawSetDomain] = React.useState(() => getInitialDomain(data));
  function setDomain(d) {
    setCrosshairValues([]);
    rawSetDomain(d);
  }
  React.useEffect(() => {
    setDomain(getInitialDomain(data));
  }, [data]);

  const items = values(data);

  return (
    <div
      className={`chart ${noSelect}`}
      style={{ flex: 1, position: 'relative' }}
    >
      <Loader loading={loading} />
      <div style={{ display: 'flex' }}></div>
      <div>
        <Plot
          animation
          height={window.innerWidth < 1000 ? 185 : 350}
          xDomain={domain && [domain.left, domain.right]}
          yDomain={domain && getYDomain(data, [domain.left, domain.right])}
          margin={{ left: 45 }}
        >
          <HorizontalGridLines />
          <VerticalGridLines />
          <XAxis
            style={{
              text: {
                fontSize: window.innerWidth < 1000 ? 7 : 11,
              },
            }}
            tickFormat={(i) => formattedDates[i]}
          />
          <YAxis tickFormat={formatNumber} />
          <ChartLabel
            text="Date"
            includeMargin={false}
            xPercent={0.025}
            yPercent={1.01}
          />

          <ChartLabel
            text={text}
            className="alt-y-label"
            includeMargin={false}
            xPercent={0.03}
            yPercent={0.06}
            style={{
              transform: 'rotate(-90)',
              textAnchor: 'end',
            }}
          />
          <Crosshair
            className={crosshairClass}
            values={crosshairValues}
            itemsFormat={(items) => items.map((i) => formatItem(formatter, i))}
            titleFormat={formatTitle}
          />
          {items.map((item, i) => (
            <LineSeries
              key={item.key}
              curve={curveCatmullRom.alpha(0.5)}
              onNearestX={
                i === 0
                  ? (value, { index }) => {
                      if (!brushing.current) {
                        setCrosshairValues(
                          items.map((i) => data[i.key].dates[index]),
                        );
                      }
                    }
                  : null
              }
              data={data[item.key].dates}
            />
          ))}
          {crosshairValues[0]
            ? items.map((item) => (
                <MarkSeries
                  animation={false}
                  key={item.key}
                  stroke="white"
                  data={[data[item.key].dates[crosshairValues[0].index]]}
                />
              ))
            : null}
        </Plot>
      </div>
      <div>
        <Brusher
          domain={domain}
          setDomain={setDomain}
          brushing={brushing}
          data={data}
          items={items}
        />
      </div>
      <DiscreteColorLegend
        onItemClick={selected.length > 0 ? onLegendClick : undefined}
        className={legendClass}
        items={items.map((i) => i.displayName)}
        orientation="horizontal"
      />
    </div>
  );
}

function isValid(data, area) {
  if (!area) return null;
  const first = dateRange;
  return area.left >= 0 && area.right <= first.length - 1;
}

const Brusher = React.memo(
  ({ setDomain: rawSetDomain, domain, brushing, data, items }) => {
    const [area, rawSetArea] = React.useState(() => getInitialDomain(data));
    const setArea = React.useMemo(
      () => (v, domain) => {
        if (!isValid(data, v)) {
          return rawSetArea(domain);
        }
        rawSetArea(v);
      },
      [rawSetArea, data],
    );
    const setDomain = React.useMemo(
      () => (v) => {
        if (!isValid(data, v)) {
          return;
        }
        rawSetDomain(v);
      },
      [rawSetDomain, data],
    );
    React.useEffect(() => {
      setArea(getInitialDomain(data));
    }, [setArea, data]);

    return (
      <Plot
        animation
        height={window.innerWidth < 1000 ? 60 : 100}
        style={{ overflow: 'visible' }}
        margin={{ left: 40, right: 10, top: 10, bottom: 0 }}
      >
        {items.map((item, i) => (
          <LineSeries
            key={item.key}
            curve={curveCatmullRom.alpha(0.5)}
            data={data[item.key].dates}
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
            setArea(area, domain);
            setDomain(area);
          }}
          onDrag={(a) => {
            setDomain(a);
          }}
          onDragEnd={(a) => {
            setArea(a, domain);
            setDomain(a);
          }}
        />
      </Plot>
    );
  },
);
