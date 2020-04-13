import React from 'react';
import '../node_modules/react-vis/dist/style.css';
import { css } from '@emotion/core';
import WorldPage from './WorldPage';
import UsPage from './UsPage';
import LeafletPage from './LeafletPage';
import Header from './Header';
import OpenLayersPage from './OpenLayersPage';

import { Router } from '@reach/router';

function App() {
  return (
    <div
      className={`App`}
      css={css`
        display: flex;
        flex-direction: column;
        flex: 1;
        min-height: 100vh;
        path {
          transition: fill 0.3s;
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
      <Router style={{ flex: 1, display: 'flex' }}>
        <WorldPage path="/old_world" />
        <OpenLayersPage path="/open_layers" />
        <UsPage path="/us" />
        <LeafletPage path="/" />
      </Router>
    </div>
  );
}

export default App;
