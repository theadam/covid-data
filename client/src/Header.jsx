import React from 'react';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import { css } from '@emotion/core';

export default function Header() {
  return (
    <AppBar
      color="primary"
      position="static"
      css={css`
        margin-bottom: 20px;
      `}
    >
      <Toolbar>
        <Typography variant="h6" color="inherit">
          Covid Data
        </Typography>
      </Toolbar>
    </AppBar>
  );
}
