import React from 'react';
import '../node_modules/react-vis/dist/style.css';
import { css } from '@emotion/core';
import WorldPage from './WorldPage';
import UsPage from './UsPage';
import Header from './Header';

import { Redirect, Router } from '@reach/router';

function App() {
  return (
    <div
      className={`App`}
      css={css`
        display: flex;
        flex-direction: column;
        flex: 1;
        margin-left: 30px;
        margin-right: 30px;
        path {
          transition: fill 0.2s;
        }
        path:hover {
          opacity: 0.5;
        }
        .chart-tip {
          display: none;
        }
        svg:hover + .chart-tip {
          display: block;
        }
      `}
    >
      <Header />
      <Router>
        <Redirect from="/" to="/world" noThrow />
        <WorldPage path="/world" />
        <UsPage path="/us" />
      </Router>
    </div>
  );
}

export default App;
