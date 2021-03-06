import React from 'react';
import Autocomplete from './AutocompleteInput';
import Typography from '@material-ui/core/Typography';
import { worldItem } from './features';
import { formatPercent, perMillionPop, formatDate, isToday } from './utils';

export default function BasicDataSection({ selectedItem, onSelect, index }) {
  const item = selectedItem ? selectedItem : worldItem;
  const date = item.dates[index];

  const confirmed = date.confirmed;
  const deaths = date.deaths;
  const deathPercent = formatPercent(deaths / confirmed);

  return (
    <div style={{ color: '#333' }}>
      <div
        style={{ marginBottom: 18, display: 'flex', justifyContent: 'center' }}
      >
        <Autocomplete
          selected={selectedItem}
          onChange={(_, item) => onSelect(item)}
        />
      </div>
      <div style={{ marginBottom: 12 }}>
        <Typography variant="h4">
          {item.displayName} Data
          <Typography variant="subtitle2" component="div">
            Population {item.population.toLocaleString()}
          </Typography>
          <Typography variant="caption" component="div">
            As of {!isToday(date.date) ? formatDate(date.date) : 'Today'}
          </Typography>
        </Typography>
      </div>
      <div style={{ marginBottom: 12 }}>
        <Typography variant="h6">
          {confirmed.toLocaleString()} Confirmed Cases
          <Typography variant="subtitle2" component="span">
            {' '}
            ({perMillionPop(confirmed, item.population)} per 1m population)
          </Typography>
        </Typography>
      </div>

      <Typography variant="h6">
        {deaths.toLocaleString()} Deaths
        <Typography variant="subtitle2" component="span">
          ({perMillionPop(deaths, item.population)} per 1m population)
          <br />
          {deathPercent.toLocaleString()}% of Confirmed Cases
        </Typography>
      </Typography>
    </div>
  );
}
