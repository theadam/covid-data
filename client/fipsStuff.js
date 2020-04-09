const json = require('./fips_counties.json');
const stateCodes = require('./src/stateCodes.json');

const result = json.reduce((acc, v) => {
  const key = v.statefp + v.countyfp;
  const state = stateCodes[v.state];
  return {
    ...acc,
    [key]: {
      name: v.countyname,
      stateCode: v.state,
      state,
      displayName: `${v.countyname}, ${state}`,
    },
  };
}, {});

console.log(JSON.stringify(result, 0, ' '));
