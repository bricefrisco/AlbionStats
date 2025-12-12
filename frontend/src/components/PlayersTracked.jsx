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

// API endpoint configuration
const API_BASE_URL = 'https://api.bricefrisco.com';
const METRICS_ENDPOINT = '/albionstats/v1/metrics/players_total';

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

// Simple min implementation
const min = (data, accessor) => {
  if (data.length === 0) return 0;

  let minimum = Infinity;
  for (const item of data) {
    const value = accessor(item);
    if (value < minimum) minimum = value;
  }

  return minimum;
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
  const axisMargin = { top: 10, right: 0, bottom: 30, left: 50 };
  const containerRef = useRef(null);
  const [dimensions, setDimensions] = useState({ width: 800, height: 400 });

  // API data state
  const [playerData, setPlayerData] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [dimensionsMeasured, setDimensionsMeasured] = useState(false);

  useEffect(() => {
    const updateDimensions = () => {
      if (containerRef.current) {
        const rect = containerRef.current.getBoundingClientRect();
        const newWidth = rect.width;
        const newHeight = rect.height;

        console.log('Container dimensions:', {
          width: newWidth,
          height: newHeight,
          container: containerRef.current,
          computedStyle: window.getComputedStyle(containerRef.current),
        });

        // Only update if we have meaningful dimensions and reasonable height
        if (newWidth > 0 && newHeight > 50) {
          // Require at least 50px height
          console.log('Setting dimensions to measured values:', {
            width: newWidth,
            height: newHeight,
          });
          setDimensions({
            width: newWidth,
            height: newHeight,
          });
          setDimensionsMeasured(true);
        } else if (newWidth > 0) {
          // If height is too small, use container height but measured width
          console.log('Height too small, falling back to 400px');
          setDimensions({
            width: newWidth,
            height: 400, // Fallback to container height
          });
          setDimensionsMeasured(true);
        }
      }
    };

    // Use setTimeout to ensure DOM has been painted
    const timeoutId = setTimeout(updateDimensions, 100);

    window.addEventListener('resize', updateDimensions);
    return () => {
      clearTimeout(timeoutId);
      window.removeEventListener('resize', updateDimensions);
    };
  }, []);

  // Fetch player data from API
  useEffect(() => {
    const fetchPlayerData = async () => {
      try {
        setLoading(true);
        const response = await fetch(
          `${API_BASE_URL}${METRICS_ENDPOINT}?granularity=1w`
        );
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        const responseData = await response.json();
        const apiData = responseData.data;

        console.log('API Response data:', apiData);

        // Transform API response to component format
        const transformedData = apiData.map((item) => ({
          date: item.timestamp, // Keep full timestamp for hourly granularity
          players: item.value,
        }));

        console.log('Transformed data:', transformedData);

        setPlayerData(transformedData);
        setError(null);
      } catch (err) {
        console.error('Failed to fetch player data:', err);
        setError('Failed to load player data');
        // Fallback to empty data
        setPlayerData([]);
      } finally {
        setLoading(false);
      }
    };

    fetchPlayerData();
  }, []);

  const width = dimensions.width;
  const height = dimensions.height;

  // Always call hooks before any conditional returns
  // scales (account for axis margins) - only create if we have data
  const dateScale = useMemo(() => {
    if (playerData.length === 0) return null;
    return scaleTime({
      range: [axisMargin.left, width - axisMargin.right],
      domain: extent(playerData, getDate),
      nice: true,
    });
  }, [width, axisMargin.left, axisMargin.right, playerData]);

  const playerScale = useMemo(() => {
    if (playerData.length === 0) return null;

    const minVal = min(playerData, getPlayerCount);
    const maxVal = max(playerData, getPlayerCount);
    const difference = maxVal - minVal;
    const yAxisRange = difference * 2; // Double the difference

    // Center the data within the range
    const center = (minVal + maxVal) / 2;

    return scaleLinear({
      range: [height - axisMargin.bottom, axisMargin.top],
      domain: [center - yAxisRange / 2, center + yAxisRange / 2],
      nice: true,
    });
  }, [height, axisMargin.top, axisMargin.bottom, playerData]);

  // tooltip handler
  const handleTooltip = useCallback(
    (event) => {
      if (!dateScale || !playerScale || playerData.length === 0) return;

      const { x } = localPoint(event) || { x: 0 };
      const x0 = dateScale.invert(x);
      const index = bisectDate(playerData, x0, 1);
      const d0 = playerData[index - 1];
      const d1 = playerData[index];

      // Guard against undefined data points (outside data range)
      if (!d0) return;

      // Snap to nearest actual data point for exact curve alignment
      let d = d0;
      if (d1 && getDate(d1)) {
        d =
          x0.getTime() - getDate(d0).getTime() >
          getDate(d1).getTime() - x0.getTime()
            ? d1
            : d0;
      }

      // Position tooltip exactly at the data point (which is on the curve)
      showTooltip({
        tooltipData: d,
        tooltipLeft: dateScale(getDate(d)),
        tooltipTop: playerScale(getPlayerCount(d)),
      });
    },
    [showTooltip, playerScale, dateScale, playerData]
  );

  if (width < 10) return null;

  // Show loading state
  if (loading) {
    return (
      <div
        ref={containerRef}
        style={{
          width: '100%',
          height: '100%',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
        }}
      >
        <div style={{ color: '#9CA3AF' }}>Loading player data...</div>
      </div>
    );
  }

  // Show error state
  if (error) {
    return (
      <div
        ref={containerRef}
        style={{
          width: '100%',
          height: '100%',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
        }}
      >
        <div style={{ color: '#EF4444' }}>{error}</div>
      </div>
    );
  }

  // Show empty state
  if (playerData.length === 0) {
    return (
      <div
        ref={containerRef}
        style={{
          width: '100%',
          height: '100%',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
        }}
      >
        <div style={{ color: '#9CA3AF' }}>No player data available</div>
      </div>
    );
  }

  // At this point, we know we have data and scales are created
  // Only render chart when dimensions are properly measured
  if (!dimensionsMeasured) {
    return (
      <div
        ref={containerRef}
        style={{
          width: '100%',
          height: '100%',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
        }}
      >
        <div style={{ color: '#9CA3AF' }}>Preparing chart...</div>
      </div>
    );
  }

  // bounds (no margins, chart fills entire area)

  return (
    <div
      ref={containerRef}
      style={{
        width: '100%',
        height: '400px',
        position: 'relative',
      }}
    >
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
          numTicks={7}
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
