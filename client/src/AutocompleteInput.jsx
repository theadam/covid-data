import React from 'react';
import Autocomplete, {
  createFilterOptions,
} from '@material-ui/lab/Autocomplete';
import TextField from '@material-ui/core/TextField';
import { allDataValues } from './features';
import InputAdornment from '@material-ui/core/InputAdornment';
import PublicIcon from '@material-ui/icons/Public';

const filterOptions = createFilterOptions({
  limit: 100,
});

export default function ({ selected, onChange }) {
  return (
    <Autocomplete
      id="combo-box-demo"
      filterOptions={filterOptions}
      options={allDataValues}
      getOptionLabel={(option) => option.displayName}
      style={{ width: 300 }}
      selectOnFocus
      renderInput={({ value, onChange, ...params }) => (
        <TextField
          {...params}
          InputProps={{
            ...params.InputProps,
            startAdornment: (
              <InputAdornment position="start">
                <PublicIcon color="primary" />
              </InputAdornment>
            ),
          }}
          label="Location Search"
          variant="outlined"
          size="small"
        />
      )}
      value={selected}
      onChange={onChange}
    />
  );
}
