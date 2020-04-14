import React from 'react';
import Snackbar from '@material-ui/core/Snackbar';

let listeners = [];
const listen = (fn) => {
  listeners.push(fn);
  return () => {
    listeners = listeners.filter((f) => f !== fn);
  };
};

const emit = (val) => {
  listeners.forEach((f) => f(val));
};

export default function MyBar() {
  const [text, setText] = React.useState(null);
  React.useEffect(() => {
    return listen(setText);
  }, []);
  return (
    <Snackbar
      anchorOrigin={{
        vertical: 'bottom',
        horizontal: 'left',
      }}
      open={text !== null}
      onClose={() => setText(null)}
      autoHideDuration={6000}
      message={text}
    />
  );
}

MyBar.open = (text) => {
  emit(text);
};
