import React from 'react';

import WorldMap from './WorldMap';
import Chart from './Chart';
import { css } from 'emotion';

const flexDirection = css`
  .chart {
    margin-left: 40px;
  }
  flex: 1;
  @media only screen and (max-width: 1000px) {
    flex-direction: column;
    .chart {
      margin-left: 0;
      margin-top: 10px;
    }
  }
`;

function makeChartData(data) {
  if (!data) return data;
  const values = Object.keys(data).map((k) => data[k]);
  const d = values.reduce((acc, v) => {
    const key = v[0].country;
    acc[key] = v;
    return acc;
  }, {});
  return d;
}

export default function WorldPage() {
  const [data, setData] = React.useState(null);
  const [chartedCountries, setChartedCountries] = React.useState([]);
  React.useEffect(() => {
    fetch('/api/data/countries/historical/')
      .then((r) => r.json())
      .then(setData);
  }, []);
  const chartData = React.useMemo(() => makeChartData(data), [data]);

  function toggleCharted(name) {
    if (chartedCountries.indexOf(name) >= 0) {
      setChartedCountries(chartedCountries.filter((n) => n !== name));
    } else {
      setChartedCountries([...chartedCountries, name]);
    }
  }

  return (
    <div className={flexDirection} style={{ display: 'flex' }}>
      <WorldMap
        loading={data === null}
        data={data}
        onDataClick={(data) => toggleCharted(data.country)}
        chartedCountries={chartedCountries}
      />
      <Chart
        loading={data === null}
        data={chartData}
        chartedCountries={chartedCountries}
        onLegendClick={toggleCharted}
      />
    </div>
  );
}
