import React from 'react';
import { geoPath } from 'd3-geo';

export const debounce = (fn, time) => {
  let cancel;
  return (...args) => {
    if (cancel) clearTimeout(cancel);
    cancel = setTimeout(() => fn(...args), time);
  };
};

export function sizeProjection(projection, width) {
  const outline = { type: 'Sphere' };

  const [[x0, y0], [x1, y1]] = geoPath(
    projection.fitWidth(width, outline),
  ).bounds(outline);

  const dy = Math.ceil(y1 - y0),
    l = Math.min(Math.ceil(x1 - x0), dy);
  projection.scale((projection.scale() * (l - 1)) / l).precision(0.2);
  return dy;
}

// Get the index for each key in a map of arrays
export function getDataIndex(data, index) {
  if (!data) return [];
  const keys = Object.keys(data);
  return keys.map((key) => {
    const list = data[key];
    return list[index !== undefined ? index : data[key].length - 1];
  });
}

export function mapBy(finals, key) {
  const result = {};
  finals.forEach((final) => {
    if (!final[key]) return;
    result[final[key]] = final;
  });
  return result;
}

export function getMax(finals, key) {
  let max = 0;

  finals.forEach((datum) => {
    const val = datum[key];
    if (val > max) {
      max = val;
    }
  });
  return max;
}

export function firstArray(data) {
  if (!data) return [];
  const firstKey = Object.keys(data)[0];
  return data[firstKey];
}

const dateRegexp = /^(\d{4})-0?(\d{1,2})-0?(\d{1,2})T/;
export function formatDate(d) {
  const [, , /*year*/ month, day] = dateRegexp.exec(d);
  return `${month}/${day}`;
}

export function usePlayer(length) {
  const [frame, rawSetFrame] = React.useState(Math.max(length - 1, 0));
  const [playing, setPlaying] = React.useState(false);
  const timeoutRef = React.useRef(null);
  React.useEffect(() => {
    rawSetFrame(Math.max(length - 1, 0));
  }, [length]);

  function setFrame(l) {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
      setPlaying(false);
    }
    rawSetFrame(l);
  }

  function play() {
    if (playing) {
      clearTimeout(timeoutRef.current);
      setPlaying(false);
      return;
    }
    let cur = frame === length - 1 ? 0 : frame;
    setPlaying(true);
    const step = () => {
      rawSetFrame(cur++);
      if (cur < length) {
        timeoutRef.current = setTimeout(step, 200);
      } else {
        setPlaying(false);
      }
    };
    step();
  }
  return { frame, setFrame, playing, play };
}

export function useProjection(projection) {
  const [width, setWidthState] = React.useState(null);
  const [path, setPath] = React.useState(() =>
    geoPath().projection(projection),
  );
  const [height, setHeight] = React.useState(null);
  const setWidth = React.useMemo(
    () => (width) => {
      setWidthState(width);
      const height = sizeProjection(projection, width);
      setPath(() => geoPath().projection(projection));
      setHeight(height);
    },
    [setWidthState, projection],
  );
  const comp = React.useRef(null);

  const ref = React.useMemo(
    () => (c) => {
      if (c === null) return;
      comp.current = c;
      const width = c.getBoundingClientRect().width;
      setWidth(width);
    },
    [setWidth],
  );

  React.useEffect(() => {
    function listener() {
      if (comp.current) {
        setWidth(comp.current.getBoundingClientRect().width);
      }
    }
    const debounced = debounce(listener, 400);
    window.addEventListener('resize', debounced);
    return () => window.removeEventListener('resize', debounced);
  }, [setWidth]);

  return { ref, path, width, height };
}

export function usePrevious(value) {
  const ref = React.useRef();
  React.useEffect(() => {
    ref.current = value;
  });
  return ref.current;
}
