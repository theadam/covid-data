import React from 'react';
import Chart from './Chart';
import '../node_modules/react-vis/dist/style.css';
import WorldMap from './WorldMap';

function App() {
  const [data, setData] = React.useState(null);
  const [chartedCountries, setChartedCountries] = React.useState([]);
  React.useEffect(() => {
    fetch('/data/countries/historical/')
      .then((r) => r.json())
      .then(setData);
  }, []);

  function toggleCharted(name) {
    if (chartedCountries.indexOf(name) >= 0) {
      setChartedCountries(chartedCountries.filter((n) => n !== name));
    } else {
      setChartedCountries([...chartedCountries, name]);
    }
  }
  return (
    <div className="App" style={{ display: 'flex' }}>
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
  );
}

export default App;
