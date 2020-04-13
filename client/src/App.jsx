import React from 'react';
import '../node_modules/react-vis/dist/style.css';
import { css } from '@emotion/core';
import LeafletPage from './LeafletPage';
import Header from './Header';

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
      <LeafletPage />
    </div>
  );
}

export default App;
