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
        @media only screen and (max-width: 1000px) {
          margin-left: 10px;
          margin-right: 10px;
        }
        path {
          transition: fill 0.3s;
        }
        path.highlight:hover {
          opacity: 0.5;
        }
        .map-tip {
          display: none;
        }
        .map-container:hover + .map-tip {
          display: block;
        }
        .rv-mouse-target {
          touch-action: pan-x;
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
