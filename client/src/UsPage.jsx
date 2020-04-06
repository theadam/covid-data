import React from 'react';

import CountyMap from './CountyMap';

export default function UsPage() {
  const [data, setData] = React.useState(null);
  React.useEffect(() => {
    Promise.all([
      fetch('/api/data/us/counties/historical/').then((r) => r.json()),
      fetch('/api/data/us/states/historical/').then((r) => r.json()),
    ]).then(([rc, rs]) => setData({ counties: rc, states: rs }));
  }, []);

  return (
    <CountyMap
      counties={data && data.counties}
      states={data && data.states}
      loading={data === null}
    />
  );
}
