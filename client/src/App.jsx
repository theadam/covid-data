import React from 'react';
import '../node_modules/react-vis/dist/style.css';
import { Global, css } from '@emotion/core';
import Header from './Header';
import styled from '@emotion/styled';
import MapCard from './MapCard';
import ChartCard from './ChartCard';

const cardWidth = 650;

const Body = styled.div`
  width: ${cardWidth}px;
  align-self: center;
  margin: auto;
  @media only screen and (max-width: ${cardWidth}px) {
    width: 100%;
  }
`;

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
      <Global
        styles={css`
          body {
            background-color: #f5f5f5;
          }
        `}
      />
      <Header />
      <Body>
        <MapCard />
        <ChartCard />
      </Body>
    </div>
  );
}

export default App;
