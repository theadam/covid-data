import React from 'react';
import IconButton from '@material-ui/core/IconButton';
import PlayArrow from '@material-ui/icons/PlayArrow';
import Pause from '@material-ui/icons/Pause';

import AirbnbSlider from './AirbnbSlider';
export default function PlaySlider({
  length,
  play,
  index,
  playing,
  setIndex,
  formatLabel,
  hideTip,
}) {
  return (
    <div style={{ display: 'flex', alignItems: 'center' }}>
      <IconButton onClick={play}>
        {playing ? <Pause /> : <PlayArrow />}
      </IconButton>
      <AirbnbSlider
        valueLabelDisplay={hideTip ? 'off' : playing ? 'on' : 'auto'}
        value={index}
        min={0}
        max={length - 1}
        marks
        onChange={(_, i) => setIndex(i)}
        valueLabelFormat={formatLabel}
      />
    </div>
  );
}
