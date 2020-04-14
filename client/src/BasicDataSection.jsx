import React from 'react';
import Autocomplete from './AutocompleteInput';
import Typography from '@material-ui/core/Typography';
import { worldItem } from './features';
import { formatDate, isToday } from './utils';

function changeStats(current, before) {
  if (!before) return null;
  const confirmedChange = current.confirmed - before.confirmed;
  const deathsChange = current.deaths - before.deaths;
  return {
    confirmedChange,
    deathsChange,
    confirmedPercent:
      Math.round((confirmedChange / current.confirmed) * 1000) / 10 || 0,
    deathsPercent: Math.round((deathsChange / current.deaths) * 1000) / 10 || 0,
  };
}

export default function BasicDataSection({ selectedItem, onSelect, index }) {
  const item = selectedItem ? selectedItem : worldItem;
  const date = item.dates[index];
  const dateBefore = item.dates[index - 1];

  const confirmed = date.confirmed;
  const deaths = date.deaths;
  const stats = changeStats(date, dateBefore);

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
          <Typography variant="caption" component="div">
            As of {!isToday(date.date) ? formatDate(date.date) : 'Today'}
          </Typography>
        </Typography>
      </div>
      <div style={{ marginBottom: 12 }}>
        <Typography variant="h6">
          {confirmed.toLocaleString()} Confirmed Cases
          {stats && (
            <Typography variant="subtitle2" component="div">
              {stats.confirmedChange.toLocaleString()} added (
              {stats.confirmedPercent}% increase)
            </Typography>
          )}
        </Typography>
      </div>

      <Typography variant="h6">
        {deaths.toLocaleString()} Deaths
        {dateBefore && (
          <Typography variant="subtitle2" component="div">
            {stats.deathsChange.toLocaleString()} added ({stats.deathsPercent}%
            increase)
          </Typography>
        )}
      </Typography>
    </div>
  );
}
