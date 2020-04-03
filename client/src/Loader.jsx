import React from 'react';
import CircularProgress from '@material-ui/core/CircularProgress';
import { css } from '@emotion/core';

export default function Loader({ loading }) {
  if (!loading) return null;
  return (
    <div
      css={css`
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translateY(-50%);
      `}
    >
      <CircularProgress />
    </div>
  );
}
