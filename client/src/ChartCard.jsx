import React from 'react';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import IconButton from '@material-ui/core/IconButton';
import MoreVertIcon from '@material-ui/icons/MoreHoriz';
import { css } from '@emotion/core';
import Autocomplete from './AutocompleteInput';
import Chart, { typeText } from './Chart';
import ClickAwayListener from '@material-ui/core/ClickAwayListener';
import Grow from '@material-ui/core/Grow';
import Paper from '@material-ui/core/Paper';
import Popper from '@material-ui/core/Popper';
import MenuItem from '@material-ui/core/MenuItem';
import MenuList from '@material-ui/core/MenuList';

export default function ChartCard() {
  const [open, setOpen] = React.useState(false);
  const [type, setType] = React.useState('confirmed');
  const [selectedItems, setSelectedItems] = React.useState([]);
  function handleLegendClick(name) {
    setSelectedItems(selectedItems.filter((item) => item.displayName !== name));
  }
  const anchor = React.useRef();

  function handleType(type) {
    return () => {
      setType(type);
      setOpen(false);
    };
  }

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
          <div
            style={{
              marginBottom: 18,
              display: 'flex',
              justifyContent: 'center',
            }}
          >
            <Autocomplete
              multiple
              selected={selectedItems}
              onChange={(_, item) => setSelectedItems(item)}
            />
          </div>
        </CardContent>
        <Chart
          selected={selectedItems}
          onLegendClick={handleLegendClick}
          type={type}
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
                  {Object.keys(typeText).map((type) => (
                    <MenuItem onClick={handleType(type)} key={type}>
                      {typeText[type]}
                    </MenuItem>
                  ))}
                </MenuList>
              </ClickAwayListener>
            </Paper>
          </Grow>
        )}
      </Popper>
    </>
  );
}
