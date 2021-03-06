import React from 'react';
import { geoPath } from 'd3-geo';
import { polyline } from 'leaflet';

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

export function mapBy(finals, fn) {
  const result = {};
  finals.forEach((final) => {
    if (!final) return;
    const key = fn(final);
    if (!key) return;
    result[key] = final;
  });
  return result;
}

export function getAllMax(data, getter) {
  let max = 0;
  const keys = Object.keys(data || {});

  keys.forEach((key) => {
    const item = data[key];
    const dates = item.dates;
    if (!dates[dates.length - 1])
      console.log(key, item, dates[dates.length - 1]);
    const val = getter(item, dates.length - 1);
    if (val > max) {
      max = val;
    }
  });
  return max;
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
  return data[firstKey].dates;
}

const dateRegexp = /^(\d{4})-0?(\d{1,2})-0?(\d{1,2})T?/;
export function formatDate(d) {
  const [, , /*year*/ month, day] = dateRegexp.exec(d);
  return `${month}/${day}`;
}

export function toDate(d) {
  const [, year, month, day] = dateRegexp.exec(d);
  return new Date(Number(year), Number(month) - 1, Number(day));
}

export function isToday(d) {
  const date = toDate(d);
  const today = new Date();
  return (
    date.getDate() === today.getDate() &&
    date.getMonth() === today.getMonth() &&
    date.getFullYear() === today.getFullYear()
  );
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
  React.useEffect(() => {
    setPath(() => geoPath().projection(projection));
  }, [projection]);
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

export function scaleXy([x, y], scale) {
  return [x * scale, y * scale];
}

export function translateXy([x, y], translate) {
  return [x + translate[0], y + translate[1]];
}

export function transformXy(xy, transform) {
  return translateXy(
    scaleXy(xy, transform.scale),
    transform.translate.map((t) => t),
  );
}

export function transformBounds(bounds, transform) {
  return bounds.map((xy) => transformXy(xy, transform));
}

export function covidTipInfo(data, showNoCases) {
  if (!data) {
    if (showNoCases) {
      return <div>No Cases</div>;
    } else {
      return null;
    }
  }
  return (
    <>
      <div>{data.confirmed} Confirmed Cases</div>
      <div>{data.deaths} Fatalities</div>
    </>
  );
}

const covidTipMaker = (showNoCases) => (feature, data) => {
  return {
    title: feature.properties.name,
    info: covidTipInfo(data, showNoCases),
  };
};

export const covidTip = covidTipMaker(false);
export const covidTipIncludingNoCases = covidTipMaker(true);

export function pluck(keys, map) {
  return keys.reduce((acc, key) => ({ ...acc, [key]: map[key] }), {});
}

export function sort(list, fn) {
  const copy = list.slice();
  copy.sort(fn);
  return copy;
}

function flipCoords(coords) {
  if (coords.length === 0) return coords;
  if (Array.isArray(coords[0])) return coords.map(flipCoords);
  return [coords[1], coords[0]];
}

export function makePolyline(geometry) {
  const coords = geometry.coordinates.map(flipCoords);
  return polyline(coords);
}

export function values(obj) {
  const keys = Object.keys(obj);
  return keys.map((key) => obj[key]);
}

export function mapObject(obj, fn) {
  return Object.keys(obj).reduce((acc, key) => {
    return { ...acc, [key]: fn(obj[key], key) };
  }, {});
}

export function roundTenths(p) {
  return Math.round(p * 10) / 10;
}

export function formatPercent(p) {
  return roundTenths(p * 100);
}

export function empty(val) {
  return val === null || val === undefined;
}

export function perMillionPop(value, pop) {
  if (empty(value) || !pop) return undefined;
  const ratio = pop / 1000000;
  return Math.round(value / ratio);
}
