import React, {
  useMemo,
  useCallback,
  useRef,
  useEffect,
  useState,
} from 'react';
import { AreaClosed, Line, Bar } from '@visx/shape';
import { curveMonotoneX } from '@visx/curve';
import { GridRows, GridColumns } from '@visx/grid';
import { AxisLeft, AxisBottom } from '@visx/axis';
import { scaleTime, scaleLinear } from '@visx/scale';
import {
  withTooltip,
  Tooltip,
  TooltipWithBounds,
  defaultStyles,
} from '@visx/tooltip';
import { localPoint } from '@visx/event';
import { LinearGradient } from '@visx/gradient';

// Reusable Time Series Chart Component
const TimeSeriesChart = withTooltip(({
  width,
  height,
  data,
  xAccessor,
  yAccessor,
  xFormatter,
  yFormatter,
  xAxisLabel,
  yAxisLabel,
  colors = {
    background: '#1a1a1a',
    background2: '#0a0a0a',
    accentColor: '#60a5fa',
    accentColorDark: '#3b82f6',
    gridColor: '#374151',
    textColor: '#9CA3AF',
  },
  margins = { top: 14, right: 0, bottom: 30, left: 54 },
  showTooltip,
  hideTooltip,
  tooltipData,
  tooltipTop = 0,
  tooltipLeft = 0,
}) => {
  const containerRef = useRef(null);
  const [dimensions, setDimensions] = useState({ width: 800, height: 400 });
  const [dimensionsMeasured, setDimensionsMeasured] = useState(false);

  useEffect(() => {
    const updateDimensions = () => {
      if (containerRef.current) {
        const rect = containerRef.current.getBoundingClientRect();
        const newWidth = rect.width;
        const newHeight = rect.height;

        if (newWidth > 0 && newHeight > 50) {
          setDimensions({
            width: newWidth,
            height: newHeight,
          });
          setDimensionsMeasured(true);
        } else if (newWidth > 0) {
          setDimensions({
            width: newWidth,
            height: 400,
          });
          setDimensionsMeasured(true);
        }
      }
    };

    const timeoutId = setTimeout(updateDimensions, 100);
    window.addEventListener('resize', updateDimensions);
    return () => {
      clearTimeout(timeoutId);
      window.removeEventListener('resize', updateDimensions);
    };
  }, []);

  const { width: containerWidth, height: containerHeight } = dimensions;

  // Always call hooks before any conditional returns
  // Scales (safely handle empty/missing data)
  const xScale = useMemo(() => {
    if (!data || data.length === 0) return null;
    return scaleTime({
      range: [margins.left, containerWidth - margins.right],
      domain: [new Date(Math.min(...data.map(xAccessor))), new Date(Math.max(...data.map(xAccessor)))],
      nice: true,
    });
  }, [containerWidth, margins.left, margins.right, data, xAccessor]);

  const yScale = useMemo(() => {
    if (!data || data.length === 0) return null;
    const minVal = Math.min(...data.map(yAccessor));
    const maxVal = Math.max(...data.map(yAccessor));
    const difference = maxVal - minVal;
    const yAxisRange = difference * 2;

    const center = (minVal + maxVal) / 2;

    return scaleLinear({
      range: [containerHeight - margins.bottom, margins.top],
      domain: [center - yAxisRange / 2, center + yAxisRange / 2],
      nice: true,
    });
  }, [containerHeight, margins.top, margins.bottom, data, yAccessor]);

  // Tooltip handler
  const handleTooltip = useCallback(
    (event) => {
      if (!xScale || !yScale || !data || data.length === 0) return;

      const { x } = localPoint(event) || { x: 0 };
      const x0 = xScale.invert(x);

      // Simple nearest neighbor for exact alignment
      let nearestPoint = data[0];
      let minDistance = Infinity;

      for (const point of data) {
        const distance = Math.abs(x0.getTime() - xAccessor(point).getTime());
        if (distance < minDistance) {
          minDistance = distance;
          nearestPoint = point;
        }
      }

      showTooltip({
        tooltipData: nearestPoint,
        tooltipLeft: xScale(xAccessor(nearestPoint)),
        tooltipTop: yScale(yAccessor(nearestPoint)),
      });
    },
    [showTooltip, xScale, yScale, data, xAccessor, yAccessor]
  );

  // Early returns after all hooks
  if (containerWidth < 10) return null;
  if (!dimensionsMeasured) {
    return (
      <div
        ref={containerRef}
        style={{
          width: '100%',
          height: '400px',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
        }}
      >
        <div style={{ color: colors.textColor }}>Preparing chart...</div>
      </div>
    );
  }
  if (!data || data.length === 0) {
    return (
      <div
        ref={containerRef}
        style={{
          width: '100%',
          height: '400px',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
        }}
      >
        <div style={{ color: colors.textColor }}>No data available</div>
      </div>
    );
  }

  const tooltipStyles = {
    ...defaultStyles,
    background: colors.background,
    border: `1px solid ${colors.gridColor}`,
    color: 'white',
    borderRadius: '8px',
  };

  // Generate stable IDs for gradients
  const backgroundGradientId = 'area-background-gradient';
  const areaGradientId = 'area-gradient';

  return (
    <div
      ref={containerRef}
      style={{
        width: '100%',
        height: '400px',
        position: 'relative',
      }}
    >
      <svg width={containerWidth} height={containerHeight} style={{ overflow: 'visible' }}>
        <rect
          x={0}
          y={0}
          width={containerWidth}
          height={containerHeight}
          fill={`url(#${backgroundGradientId})`}
          stroke="none"
          rx={12}
          ry={12}
        />
        <LinearGradient
          id={backgroundGradientId}
          from={colors.background}
          to={colors.background}
          toOpacity={0.2}
        />
        <LinearGradient
          id={areaGradientId}
          from={colors.accentColor}
          to={colors.accentColor}
          toOpacity={0.05}
        />

        <GridRows
          left={margins.left}
          scale={yScale}
          width={containerWidth - margins.left - margins.right}
          strokeDasharray="1,3"
          stroke={colors.gridColor}
          strokeOpacity={0.3}
          pointerEvents="none"
        />
        <GridColumns
          top={margins.top}
          scale={xScale}
          height={containerHeight - margins.top - margins.bottom}
          strokeDasharray="1,3"
          stroke={colors.gridColor}
          strokeOpacity={0.2}
          pointerEvents="none"
        />

        <AreaClosed
          data={data}
          x={(d) => xScale(xAccessor(d)) ?? 0}
          y={(d) => yScale(yAccessor(d)) ?? 0}
          yScale={yScale}
          strokeWidth={2}
          stroke={colors.accentColor}
          fill={`url(#${areaGradientId})`}
          curve={curveMonotoneX}
        />

        <Bar
          x={margins.left}
          y={margins.top}
          width={containerWidth - margins.left - margins.right}
          height={containerHeight - margins.top - margins.bottom}
          fill="transparent"
          onMouseMove={handleTooltip}
          onMouseLeave={() => hideTooltip()}
        />

        <AxisLeft
          left={margins.left}
          scale={yScale}
          numTicks={5}
          stroke={colors.gridColor}
          tickStroke={colors.gridColor}
          tickLabelProps={{
            fill: colors.textColor,
            fontSize: 11,
            textAnchor: 'end',
            dx: '-0.25em',
            dy: '0.25em',
          }}
          tickFormat={yFormatter}
        />

        <AxisBottom
          top={containerHeight - margins.bottom}
          scale={xScale}
          numTicks={4}
          stroke={colors.gridColor}
          tickStroke={colors.gridColor}
          tickLabelProps={{
            fill: colors.textColor,
            fontSize: 11,
            textAnchor: 'middle',
            dy: '0.25em',
          }}
          tickFormat={xFormatter}
        />

        {tooltipData && (
          <g>
            <Line
              from={{ x: tooltipLeft, y: margins.top }}
              to={{ x: tooltipLeft, y: containerHeight - margins.bottom }}
              stroke={colors.accentColorDark}
              strokeWidth={2}
              pointerEvents="none"
              strokeDasharray="5,2"
            />
            <circle
              cx={tooltipLeft}
              cy={tooltipTop + 1}
              r={4}
              fill="black"
              fillOpacity={0.1}
              stroke="black"
              strokeOpacity={0.1}
              strokeWidth={2}
              pointerEvents="none"
            />
            <circle
              cx={tooltipLeft}
              cy={tooltipTop}
              r={4}
              fill={colors.accentColorDark}
              stroke="white"
              strokeWidth={2}
              pointerEvents="none"
            />
          </g>
        )}
      </svg>

      {tooltipData && (
        <div>
          <TooltipWithBounds
            key={Math.random()}
            top={tooltipTop - 12}
            left={tooltipLeft + 12}
            style={tooltipStyles}
          >
            {yFormatter ? yFormatter(yAccessor(tooltipData)) : yAccessor(tooltipData)}
          </TooltipWithBounds>
          <Tooltip
            top={containerHeight - margins.bottom + 16}
            left={tooltipLeft}
            style={{
              ...defaultStyles,
              minWidth: 72,
              textAlign: 'center',
              transform: 'translateX(-50%)',
              background: colors.background,
              border: `1px solid ${colors.gridColor}`,
              color: 'white',
              borderRadius: '8px',
            }}
          >
            {(() => {
              const xValue = xAccessor(tooltipData);
              const formattedDate = xFormatter ? xFormatter(xValue) : xValue;
              const localizedTime = new Date(xValue).toLocaleTimeString(undefined, {
                hour: 'numeric',
                minute: '2-digit',
              });
              return (
                <>
                  <div>{formattedDate}</div>
                  <div style={{ fontSize: 11, opacity: 0.8 }}>{localizedTime}</div>
                </>
              );
            })()}
          </Tooltip>
        </div>
      )}
    </div>
  );
});

export default TimeSeriesChart;

