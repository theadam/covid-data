import countyData from './data/counties-10m.json';
import * as topojson from 'topojson-client';
import worldData from './data/countries-110m.json';
import provinceData from './data/canadaprovtopo.json';
import australiaData from './data/au-states.json';
import chinaData from './data/china-provinces.json';
import fipsData from './fipsData.json';

import countries from './data/world.json';
import states from './data/state.json';
import provinces from './data/province.json';
import counties from './data/county.json';
import dateRange from './data/dateRange.json';
import { values, sort, mapBy, mapObject } from './utils';

export { dateRange };

export const USA = '840';
export const Canada = '124';
export const China = '156';
export const Australia = '036';

export const countriesWithRegions = [USA, Canada, China, Australia];

const baseFeatures = {
  world: topojson.feature(worldData, worldData.objects.countries),
  usStates: topojson.feature(countyData, countyData.objects.states),
  usCounties: topojson.feature(countyData, countyData.objects.counties),
  china: topojson.feature(chinaData, chinaData.objects.CHN_adm1),
  australia: topojson.feature(australiaData, australiaData.objects.states),
  canada: topojson.feature(provinceData, provinceData.objects.canadaprov),
};

function enrichFeatures(
  featureCollection,
  featureKey = (feature) => feature.id,
  displayName = (feature) => feature.properties.name,
) {
  return {
    ...featureCollection,
    features: featureCollection.features.map((feature) => {
      const key = featureKey(feature);
      return {
        ...feature,
        key,
        displayName: displayName(feature),
      };
    }),
  };
}

function createFeatures() {
  return {
    world: enrichFeatures(baseFeatures.world, undefined, (feature) =>
      feature.properties.name === 'Georgia'
        ? 'Georgia (country)'
        : feature.properties.name,
    ),
    usStates: enrichFeatures(baseFeatures.usStates),
    usCounties: enrichFeatures(
      baseFeatures.usCounties,
      undefined,
      (feature) => {
        if (!fipsData[feature.id]) {
          console.log(feature.id);
        }
        return fipsData[feature.id].displayName;
      },
    ),
    china: enrichFeatures(
      baseFeatures.china,
      (feature) => `${China}-${feature.properties.NAME_1}`,
      (feature) => feature.properties.NAME_1,
    ),
    australia: enrichFeatures(
      baseFeatures.australia,
      (feature) => `${Australia}-${feature.properties.STATE_NAME}`,
      (feature) => feature.properties.STATE_NAME,
    ),
    canada: enrichFeatures(
      baseFeatures.canada,
      (feature) => `${Canada}-${feature.properties.name}`,
    ),
  };
}

const features = createFeatures();
export default features;

function mapCountryData(featureCollection, data) {
  const featureMap = mapBy(
    featureCollection.features,
    (feature) => feature.key,
  );
  return mapObject(data, (item, key) => {
    const feature = featureMap[key];

    return {
      ...data[key],
      key: feature ? feature.key : item.countryCode,
      displayName: feature ? feature.displayName : item.country,
      geometry: feature ? feature.geometry : null,
    };
  });
}

function pluckFeatureData(featureCollection, data) {
  return featureCollection.features.reduce(
    (acc, { geometry, displayName, key }) => {
      if (data[key]) {
        acc[key] = {
          ...data[key],
          key,
          displayName,
          geometry,
        };
      }
      return acc;
    },
    {},
  );
}

function splitData(data) {
  return {
    world: mapCountryData(features.world, data.countries),
    usStates: pluckFeatureData(features.usStates, data.states),
    usCounties: pluckFeatureData(features.usCounties, data.counties),
    china: pluckFeatureData(features.china, data.provinces),
    australia: pluckFeatureData(features.australia, data.provinces),
    canada: pluckFeatureData(features.canada, data.provinces),
  };
}
export const data = splitData({ countries, states, provinces, counties });

function mergeAll(obj) {
  const keys = Object.keys(obj);
  return keys.reduce((acc, key) => ({ ...acc, ...obj[key] }), {});
}

export const allData = mergeAll(data);

const worldValues = values(data.world);

const sortData = (data) =>
  sort(data, (a, b) => a.displayName.localeCompare(b.displayName));

export const allDataValues = [
  ...sortData(values(data.world)),
  ...sortData(values(data.usStates)),
  ...sortData(values(data.china)),
  ...sortData(values(data.australia)),
  ...sortData(values(data.canada)),
  ...sortData(values(data.usCounties)),
];

export const worldItem = {
  displayName: 'Worldwide',
  key: 'world',
  dates: dateRange.map((date, i) => {
    const forIndex = worldValues.reduce(
      (acc, item) => {
        const data = item.dates[i];
        return {
          confirmed: acc.confirmed + data.confirmed,
          deaths: acc.deaths + data.deaths,
        };
      },
      { confirmed: 0, deaths: 0 },
    );
    return { date, ...forIndex };
  }),
};
