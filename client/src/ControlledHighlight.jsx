import { Highlight } from 'react-vis';
import { getAttributeScale } from 'react-vis/dist/utils/scales-utils';

export default class ControlledHighlight extends Highlight {
  componentDidMount() {
    this.updateArea(this.props.area);
  }

  componentDidUpdate() {
    if (!this.state.brushing && !this.state.dragging) {
      this.updateArea(this.props.area);
    }
  }

  updateArea(area) {
    let newLeft = 0;
    let newRight = 0;
    if (area) {
      const { left, right } = area;
      const xScale = getAttributeScale(this.props, 'x');

      const { marginLeft } = this.props;

      newLeft = xScale(left) + marginLeft;
      newRight = xScale(right) + marginLeft;
    }

    const { left, right } = this.state.brushArea;
    if (left !== newLeft || right !== newRight) {
      const { innerHeight, marginBottom, marginTop } = this.props;
      const plotHeight = innerHeight + marginTop + marginBottom;

      this.setState({
        brushArea: {
          left: newLeft,
          right: newRight,
          top: 0,
          bottom: plotHeight,
        },
        dragArea: this.props.drag
          ? {
              left: newLeft,
              right: newRight,
              top: 0,
              bottom: plotHeight,
            }
          : null,
      });
    }
  }

  render() {
    const result = super.render();
    const rect = result.props.children[0];
    const newRect = {
      ...rect,
      props: {
        ...rect.props,
        onTouchStart: (e) => {
          const evt = {
            ...e,
            nativeEvent: {
              ...e.nativeEvent,
              type: 'touchstart',
              pageX: e.nativeEvent.touches[0].pageX,
              pageY: e.nativeEvent.touches[0].pageY,
            },
          };
          e.preventDefault();
          this.startBrushing(evt);
        },
        onTouchMove: (e) => {
          const evt = {
            ...e,
            nativeEvent: {
              ...e.nativeEvent,
              offsetX: e.nativeEvent.touches[0].pageX,
              offsetY: e.nativeEvent.touches[0].pageY,
            },
          };
          e.preventDefault();
          this.onBrush(evt);
        },
      },
    };
    const newResult = {
      ...result,
      props: {
        ...result.props,
        children: [newRect, ...result.props.children.slice(1)],
      },
    };
    return newResult;
  }
}
