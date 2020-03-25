import React from 'react';
import Chart from './Chart';
import '../node_modules/react-vis/dist/style.css';
import WorldMap from './WorldMap';

function App() {
  const [data, setData] = React.useState(null)
  const [chartedCountries, setChartedCountries] = React.useState([])
  React.useEffect(() => {
    fetch('/data/countries/historical/').then(r => r.json()).then(setData)
  }, [])
  return (
    <div className="App" style={{ display: 'flex'}}>
      {data !== null &&
        <>
          <WorldMap data={data} onDataClick={(data) => {
            if (chartedCountries.indexOf(data.country) >= 0) {
              setChartedCountries(chartedCountries.filter(name => name !== data.country));
            } else {
              setChartedCountries([...chartedCountries, data.country]);
            }
          }} chartedCountries={chartedCountries}/>
          <Chart data={data} chartedCountries={chartedCountries} />
        </>
      }
    </div>
  );
}

export default App;
