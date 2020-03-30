import React from 'react';
import Chart from './Chart';
import '../node_modules/react-vis/dist/style.css';
import WorldMap from './WorldMap';
import CountyMap from './CountyMap';
import { css } from 'emotion';

const transitionPaths = css`
  path {
    transition: fill 0.2s;
  }
  path:hover {
    opacity: 0.5;
  }
  .chart-tip {
    display: none;
  }
  svg:hover + .chart-tip {
    display: block;
  }
`;

function App() {
  const [data, setData] = React.useState(null);
  const [chartedCountries, setChartedCountries] = React.useState([]);
  React.useEffect(() => {
    fetch('/api/data/countries/historical/')
      .then((r) => r.json())
      .then(setData);
  }, []);

  const [countyData, setCountyData] = React.useState(null);
  React.useEffect(() => {
    fetch('/api/data/us/counties/historical/')
      .then((r) => r.json())
      .then(setCountyData);
  }, []);

  function toggleCharted(name) {
    if (chartedCountries.indexOf(name) >= 0) {
      setChartedCountries(chartedCountries.filter((n) => n !== name));
    } else {
      setChartedCountries([...chartedCountries, name]);
    }
  }
  return (
    <div
      className={`App ${transitionPaths}`}
      style={{ display: 'flex', flexDirection: 'column' }}
    >
      <div style={{ display: 'flex' }}>
        {data !== null && (
          <>
            <WorldMap
              data={data}
              onDataClick={(data) => toggleCharted(data.country)}
              chartedCountries={chartedCountries}
            />
            <Chart
              data={data}
              chartedCountries={chartedCountries}
              onLegendClick={toggleCharted}
            />
          </>
        )}
      </div>
      <div>{countyData !== null && <CountyMap data={countyData} />}</div>
    </div>
  );
}

export default App;
