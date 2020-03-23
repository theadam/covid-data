import React from 'react';
import { curveCatmullRom } from 'd3-shape';
import '../node_modules/react-vis/dist/style.css';

import {
  XYPlot,
  XAxis,
  YAxis,
  ChartLabel,
  HorizontalGridLines,
  VerticalGridLines,
  LineSeries,
  DiscreteColorLegend,
  makeWidthFlexible,
} from 'react-vis';
import { css } from 'emotion';

export default class Example extends React.Component {
  render() {
    const Plot = makeWidthFlexible(XYPlot);

    return (
      <div>
        <Plot height={500}>
          <HorizontalGridLines />
          <VerticalGridLines />
          <XAxis />
          <YAxis />
          <ChartLabel
            text="X Axis"
            includeMargin={false}
            xPercent={0.025}
            yPercent={1.01}
          />

          <ChartLabel
            text="Y Axis"
            className="alt-y-label"
            includeMargin={false}
            xPercent={0.01}
            yPercent={0.06}
            style={{
              transform: 'rotate(-90)',
              textAnchor: 'end',
            }}
          />
          <LineSeries
            data={[
              { x: 1, y: 3 },
              { x: 2, y: 5 },
              { x: 3, y: 15 },
              { x: 4, y: 12 },
            ]}
          />
          <LineSeries className="second-series" data={null} />
          <LineSeries
            curve={'curveMonotoneX'}
            data={[
              { x: 1, y: 10 },
              { x: 2, y: 4 },
              { x: 3, y: 2 },
              { x: 4, y: 15 },
            ]}
            strokeDasharray={'7, 3'}
          />
          <LineSeries
            curve={curveCatmullRom.alpha(0.5)}
            style={{
              // note that this can not be translated to the canvas version
              strokeDasharray: '2 2',
            }}
            data={[
              { x: 1, y: 7 },
              { x: 2, y: 11 },
              { x: 3, y: 9 },
              { x: 4, y: 2 },
            ]}
          />
        </Plot>
        <DiscreteColorLegend
          className={css`
            text-align: right;
          `}
          items={['one', 'two', 'three', 'four']}
          orientation="horizontal"
        />
      </div>
    );
  }
}
