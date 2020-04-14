import React from 'react';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import { dateRange } from './features';
import { formatDate, usePlayer } from './utils';
import PlaySlider from './PlaySlider';
import LeafletSection from './LeafletSection';
import BasicDataSection from './BasicDataSection';
import { css } from '@emotion/core';

export default function MapCard() {
  const [selectedItem, setSelectedItem] = React.useState(null);
  const { play, playing, frame: index, setFrame: setIndex } = usePlayer(
    dateRange.length,
  );

  return (
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
  );
}
