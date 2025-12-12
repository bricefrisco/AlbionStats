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

// Mock data for players tracked over time
const playerData = [
  { date: '2024-01-01', players: 8500 },
  { date: '2024-02-01', players: 9200 },
  { date: '2024-03-01', players: 10100 },
  { date: '2024-04-01', players: 11800 },
  { date: '2024-05-01', players: 12400 },
  { date: '2024-06-01', players: 12800 },
  { date: '2024-07-01', players: 13200 },
  { date: '2024-08-01', players: 13800 },
  { date: '2024-09-01', players: 14100 },
  { date: '2024-10-01', players: 14500 },
  { date: '2024-11-01', players: 14800 },
  { date: '2024-12-01', players: 15200 },
];

const background = '#1a1a1a';
const background2 = '#0a0a0a';
const accentColor = '#60a5fa';
const accentColorDark = '#3b82f6';

const tooltipStyles = {
  ...defaultStyles,
  background: background,
  border: '1px solid #374151',
  color: 'white',
  borderRadius: '8px',
};

// Simple date formatter
const formatDate = (date) => {
  return new Date(date).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: '2-digit',
  });
};

// accessors
const getDate = (d) => new Date(d.date);
const getPlayerCount = (d) => d.players;

// Simple bisector implementation
const bisectDate = (data, target) => {
  let left = 0;
  let right = data.length - 1;

  while (left <= right) {
    const mid = Math.floor((left + right) / 2);
    const midDate = getDate(data[mid]).getTime();

    if (midDate < target.getTime()) {
      left = mid + 1;
    } else {
      right = mid - 1;
    }
  }

  return left;
};

// Simple extent implementation
const extent = (data, accessor) => {
  if (data.length === 0) return [0, 0];

  let min = Infinity;
  let max = -Infinity;

  for (const item of data) {
    const value = accessor(item);
    if (value < min) min = value;
    if (value > max) max = value;
  }

  return [min, max];
};

// Simple max implementation
const max = (data, accessor) => {
  if (data.length === 0) return 0;

  let maximum = -Infinity;
  for (const item of data) {
    const value = accessor(item);
    if (value > maximum) maximum = value;
  }

  return maximum;
};

// PlayersTrackedChart component for Fast Refresh compatibility
const PlayersTrackedChart = ({
  width: widthProp,
  height: heightProp,
  showTooltip,
  hideTooltip,
  tooltipData,
  tooltipTop = 0,
  tooltipLeft = 0,
}) => {
  // Small margins for axis labels
  const axisMargin = { top: 10, right: 10, bottom: 30, left: 50 };
  const containerRef = useRef(null);
  const [dimensions, setDimensions] = useState({ width: 350, height: 280 });

  useEffect(() => {
    const updateDimensions = () => {
      if (containerRef.current) {
        const rect = containerRef.current.getBoundingClientRect();
        setDimensions({
          width: rect.width,
          height: rect.height,
        });
      }
    };

    updateDimensions();
    window.addEventListener('resize', updateDimensions);
    return () => window.removeEventListener('resize', updateDimensions);
  }, []);

  const width = dimensions.width;
  const height = dimensions.height;

  if (width < 10) return null;

  // bounds (no margins, chart fills entire area)

  // scales (account for axis margins)
  const dateScale = useMemo(
    () =>
      scaleTime({
        range: [axisMargin.left, width - axisMargin.right],
        domain: extent(playerData, getDate),
      }),
    [width, axisMargin.left, axisMargin.right]
  );

  const playerScale = useMemo(
    () =>
      scaleLinear({
        range: [height - axisMargin.bottom, axisMargin.top],
        domain: [0, max(playerData, getPlayerCount) * 1.1],
        nice: true,
      }),
    [height, axisMargin.top, axisMargin.bottom]
  );

  // tooltip handler
  const handleTooltip = useCallback(
    (event) => {
      const { x } = localPoint(event) || { x: 0 };
      const x0 = dateScale.invert(x);
      const index = bisectDate(playerData, x0);
      const d0 = playerData[index - 1];
      const d1 = playerData[index];
      let d = d0;

      if (d1 && getDate(d1)) {
        d =
          x0.getTime() - getDate(d0).getTime() >
          getDate(d1).getTime() - x0.getTime()
            ? d1
            : d0;
      }

      showTooltip({
        tooltipData: d,
        tooltipLeft: x,
        tooltipTop: playerScale(getPlayerCount(d)),
      });
    },
    [showTooltip, playerScale, dateScale]
  );

  return (
    <div ref={containerRef} style={{ width: '100%', height: '100%' }}>
      <svg width={width} height={height}>
        <rect
          x={0}
          y={0}
          width={width}
          height={height}
          fill="url(#area-background-gradient)"
          stroke="none"
        />
        <LinearGradient
          id="area-background-gradient"
          from={background}
          to={background2}
        />
        <LinearGradient
          id="area-gradient"
          from={accentColor}
          to={accentColor}
          toOpacity={0.1}
        />

        <GridRows
          left={axisMargin.left}
          scale={playerScale}
          width={width - axisMargin.left - axisMargin.right}
          strokeDasharray="1,3"
          stroke="#374151"
          strokeOpacity={0.3}
          pointerEvents="none"
        />
        <GridColumns
          top={axisMargin.top}
          scale={dateScale}
          height={height - axisMargin.top - axisMargin.bottom}
          strokeDasharray="1,3"
          stroke="#374151"
          strokeOpacity={0.3}
          pointerEvents="none"
        />

        <AreaClosed
          data={playerData}
          x={(d) => dateScale(getDate(d)) ?? 0}
          y={(d) => playerScale(getPlayerCount(d)) ?? 0}
          yScale={playerScale}
          strokeWidth={2}
          stroke={accentColor}
          fill="url(#area-gradient)"
          curve={curveMonotoneX}
        />

        <Bar
          x={0}
          y={0}
          width={width}
          height={height}
          fill="transparent"
          onMouseMove={handleTooltip}
          onMouseLeave={() => hideTooltip()}
        />

        <AxisLeft
          left={axisMargin.left}
          scale={playerScale}
          numTicks={5}
          stroke="#374151"
          tickStroke="#374151"
          tickLabelProps={{
            fill: '#9CA3AF',
            fontSize: 11,
            textAnchor: 'end',
            dx: '-0.25em',
            dy: '0.25em',
          }}
          tickFormat={(value) => `${value.toLocaleString()}`}
        />

        <AxisBottom
          top={height - axisMargin.bottom}
          scale={dateScale}
          numTicks={4}
          stroke="#374151"
          tickStroke="#374151"
          tickLabelProps={{
            fill: '#9CA3AF',
            fontSize: 11,
            textAnchor: 'middle',
            dy: '0.25em',
          }}
          tickFormat={formatDate}
        />

        {tooltipData && (
          <g>
            <Line
              from={{ x: tooltipLeft, y: axisMargin.top }}
              to={{ x: tooltipLeft, y: height - axisMargin.bottom }}
              stroke={accentColorDark}
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
              fill={accentColorDark}
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
            {`${getPlayerCount(tooltipData).toLocaleString()} players`}
          </TooltipWithBounds>
          <Tooltip
            top={height - axisMargin.bottom + 16}
            left={tooltipLeft}
            style={{
              ...defaultStyles,
              minWidth: 72,
              textAlign: 'center',
              transform: 'translateX(-50%)',
              background: background,
              border: '1px solid #374151',
              color: 'white',
              borderRadius: '8px',
            }}
          >
            {formatDate(getDate(tooltipData))}
          </Tooltip>
        </div>
      )}
    </div>
  );
};
const PlayersTrackedChartWithTooltip = withTooltip(PlayersTrackedChart);

export default PlayersTrackedChartWithTooltip;
