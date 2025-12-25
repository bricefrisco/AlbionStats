import React, {
  useMemo,
  useCallback,
  useRef,
  useEffect,
  useState,
  useId,
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
import { formatCompactNumber } from '../utils/formatters';

const formatYAxisValue = (value) => {
  if (value === undefined || value === null) return '';
  return formatCompactNumber(value);
};

const formatTooltipValue = (value) => {
  if (value === undefined || value === null) return '';

  const numericValue =
    typeof value === 'number' ? value : Number(value);

  if (!Number.isFinite(numericValue)) {
    return `${value}`;
  }

  const formatter = new Intl.NumberFormat(undefined, {
    minimumFractionDigits: Number.isInteger(numericValue) ? 0 : 2,
    maximumFractionDigits: 2,
  });

  return formatter.format(numericValue);
};

const TimeSeriesChart = withTooltip(
  ({
    data,
    xAccessor,
    yAccessor,
    xFormatter,
    yFormatter,
    colors = {
      background: '#1a1a1a',
      background2: '#0a0a0a',
      accentColor: '#60a5fa',
      accentColorDark: '#3b82f6',
      gridColor: '#374151',
      textColor: '#9CA3AF',
    },
    subtleGradient = false,
    margins = { top: 14, right: 0, bottom: 30, left: 54 },
    showTooltip,
    hideTooltip,
    tooltipData,
    tooltipTop = 0,
    tooltipLeft = 0,
  }) => {
  const containerRef = useRef(null);
  const [dimensions, setDimensions] = useState({ width: 0, height: 400 });

  useEffect(() => {
    if (!containerRef.current) return;

    const resizeObserver = new ResizeObserver((entries) => {
      for (const entry of entries) {
        const { width, height } = entry.contentRect;
        if (width > 0) {
          setDimensions({
            width,
            height: height > 50 ? height : 400,
          });
        }
      }
    });

    resizeObserver.observe(containerRef.current);
    return () => resizeObserver.disconnect();
  }, []);

    const { width: containerWidth, height: containerHeight } = dimensions;

    // Always call hooks before any conditional returns
    // Scales (safely handle empty/missing data)
    const xScale = useMemo(() => {
      if (!data || data.length === 0) return null;
    return scaleTime({
      range: [margins.left, containerWidth - margins.right],
      domain: [
        new Date(Math.min(...data.map(xAccessor))),
        new Date(Math.max(...data.map(xAccessor))),
      ],
    });
    }, [containerWidth, margins.left, margins.right, data, xAccessor]);

    const yScale = useMemo(() => {
      if (!data || data.length === 0) return null;
      const minVal = Math.min(...data.map(yAccessor));
      const maxVal = Math.max(...data.map(yAccessor));
      const difference = maxVal - minVal;
      const padding = difference * 0.25;

      const domainMin = Math.max(0, minVal - padding);
      const domainMax = maxVal + padding;

      return scaleLinear({
        range: [containerHeight - margins.bottom, margins.top],
        domain: [domainMin, domainMax],
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

    // Generate stable IDs for gradients (must be after all hooks but before early returns)
    const id = useId();
    const backgroundGradientId = `bg-${id}`;
    const areaGradientId = `area-${id}`;
    const areaFromOpacity = subtleGradient ? 0.35 : 1;
    const areaToOpacity = subtleGradient ? 0.0125 : 0.05;

    // Early returns after all hooks
    if (containerWidth < 10) {
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

    return (
      <div
        ref={containerRef}
        style={{
          width: '100%',
          height: '400px',
          position: 'relative',
        }}
      >
        <svg
          width={containerWidth}
          height={containerHeight}
          style={{ overflow: 'visible' }}
        >
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
            fromOpacity={areaFromOpacity}
            to={colors.accentColor}
            toOpacity={areaToOpacity}
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
          tickFormat={yFormatter ?? formatYAxisValue}
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
              {(() => {
                const yValue = yAccessor(tooltipData);
                return yFormatter
                  ? yFormatter(yValue)
                  : formatTooltipValue(yValue);
              })()}
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
                const localizedTime = new Date(xValue).toLocaleTimeString(
                  undefined,
                  {
                    hour: 'numeric',
                    minute: '2-digit',
                  }
                );
                return (
                  <>
                    <div>{formattedDate}</div>
                    <div style={{ fontSize: 11, opacity: 0.8 }}>
                      {localizedTime}
                    </div>
                  </>
                );
              })()}
            </Tooltip>
          </div>
        )}
      </div>
    );
  }
);

export default TimeSeriesChart;
