import React from 'react';
import '../node_modules/react-vis/dist/style.css';
import { Global, css } from '@emotion/core';
import LeafletSection from './LeafletSection';
import Header from './Header';
import BasicDataSection from './BasicDataSection';
import styled from '@emotion/styled';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import { dateRange } from './features';
import { formatDate, usePlayer } from './utils';
import PlaySlider from './PlaySlider';

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
  const [selectedItem, setSelectedItem] = React.useState(null);
  const { play, playing, frame: index, setFrame: setIndex } = usePlayer(
    dateRange.length,
  );

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
        <Card
          css={css`
            margin-bottom: 20px;
          `}
        >
          <CardContent>
            <BasicDataSection
              index={index}
              onSelect={setSelectedItem}
              selectedItem={selectedItem}
            />
          </CardContent>
          <PlaySlider
            playing={playing}
            play={play}
            index={index}
            length={dateRange.length}
            setIndex={setIndex}
            formatLabel={(i) => formatDate(dateRange[i])}
            hideTip={false}
          />
          <LeafletSection
            index={index}
            onSelect={setSelectedItem}
            centeredItem={selectedItem}
          />
        </Card>
      </Body>
    </div>
  );
}

export default App;
