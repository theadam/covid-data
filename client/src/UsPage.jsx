import React from 'react';

import CountyMap from './CountyMap';

export default function UsPage() {
  const [data, setData] = React.useState(null);
  React.useEffect(() => {
    fetch('/api/data/us/counties/historical/')
      .then((r) => r.json())
      .then(setData);
  }, []);

  return <CountyMap data={data} loading={data === null} />;
}
