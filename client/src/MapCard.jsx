import React from 'react';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import { dateRange } from './features';
import { formatDate, usePlayer, perMillionPop } from './utils';
import PlaySlider from './PlaySlider';
import LeafletSection from './LeafletSection';
import BasicDataSection from './BasicDataSection';
import { css } from '@emotion/core';
import IconButton from '@material-ui/core/IconButton';
import ClickAwayListener from '@material-ui/core/ClickAwayListener';
import MoreVertIcon from '@material-ui/icons/MoreHoriz';
import Grow from '@material-ui/core/Grow';
import Paper from '@material-ui/core/Paper';
import Popper from '@material-ui/core/Popper';
import MenuItem from '@material-ui/core/MenuItem';
import MenuList from '@material-ui/core/MenuList';

export default function MapCard() {
  const [open, setOpen] = React.useState(false);
  const [selectedItem, setSelectedItem] = React.useState(null);
  const { play, playing, frame: index, setFrame: setIndex } = usePlayer(
    dateRange.length,
  );
  const [dataKey, setDataKey] = React.useState('confirmed');
  const [usePopulation, setUsePopulation] = React.useState(false);
  const anchor = React.useRef();

  const handleMenu = React.useCallback((key, usePop = false) => {
    return () => {
      setDataKey(key);
      setUsePopulation(usePop);
      setOpen(false);
    };
  }, []);

  const calculateValue = React.useCallback(
    (item, index) => {
      if (usePopulation) {
        return perMillionPop(item?.dates?.[index]?.[dataKey], item?.population);
      }
      return item?.dates?.[index]?.[dataKey];
    },
    [usePopulation, dataKey],
  );

  return (
    <>
      <Card
        css={css`
          margin-bottom: 20px;
          position: relative;
        `}
      >
        <CardContent>
          <div
            ref={anchor}
            style={{
              position: 'absolute',
              top: 10,
              right: 10,
            }}
          >
            <IconButton aria-label="settings" onClick={() => setOpen(!open)}>
              <MoreVertIcon />
            </IconButton>
          </div>
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
          calculateValue={calculateValue}
          index={index}
          onSelect={setSelectedItem}
          centeredItem={selectedItem}
        />
      </Card>
      <Popper
        open={open}
        anchorEl={anchor.current}
        role={undefined}
        transition
        disablePortal
      >
        {({ TransitionProps, placement }) => (
          <Grow
            {...TransitionProps}
            style={{
              transformOrigin:
                placement === 'bottom' ? 'center top' : 'center bottom',
            }}
          >
            <Paper>
              <ClickAwayListener onClickAway={() => setOpen(false)}>
                <MenuList autoFocusItem={open} id="menu-list-grow">
                  <MenuItem onClick={handleMenu('confirmed')}>
                    Confirmed Cases
                  </MenuItem>
                  <MenuItem onClick={handleMenu('deaths')}>Deaths</MenuItem>
                  <MenuItem onClick={handleMenu('confirmed', true)}>
                    Confirmed Cases Per 1m Population
                  </MenuItem>
                  <MenuItem onClick={handleMenu('deaths', true)}>
                    Deaths Per 1m Population
                  </MenuItem>
                </MenuList>
              </ClickAwayListener>
            </Paper>
          </Grow>
        )}
      </Popper>
    </>
  );
}
