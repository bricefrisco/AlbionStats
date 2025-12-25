import React, { useEffect, useState } from 'react';
import { withTooltip } from '@visx/tooltip';
import TimeSeriesChart from './TimeSeriesChart';
import TimeRangeToggle from './TimeRangeToggle';
import { formatDate } from '../utils/formatters';

// API endpoint configuration
const API_BASE_URL = 'https://api.bricefrisco.com';
const METRICS_ENDPOINT = '/albionstats/v1/metrics/snapshots';

// accessors
const getDate = (d) => new Date(d.date);
const getDataPoints = (d) => d.datapoints;

// Green color scheme for the chart
const greenColors = {
  background: '#1a1a1a',
  background2: '#0a0a0a',
  accentColor: '#15803d',
  accentColorDark: '#166534',
  gridColor: '#374151',
  textColor: '#9CA3AF',
};

// DataPointsChart component for Fast Refresh compatibility
const DataPointsChart = ({
  width: widthProp,
  height: heightProp,
  showTooltip,
  hideTooltip,
  tooltipData,
  tooltipTop = 0,
  tooltipLeft = 0,
  timeRange = '1w',
}) => {
  // API data state
  const [dataPointsData, setDataPointsData] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Fetch data points data from API
  useEffect(() => {
    const fetchDataPointsData = async () => {
      try {
        setLoading(true);
        const response = await fetch(
          `${API_BASE_URL}${METRICS_ENDPOINT}?granularity=${timeRange}`
        );
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        const responseData = await response.json();
        const apiData = responseData.data;

        // Transform API response to component format
        const transformedData = apiData.map((item) => ({
          date: item.timestamp, // Keep full timestamp for hourly granularity
          datapoints: item.value,
        }));

        setDataPointsData(transformedData);
        setError(null);
      } catch (err) {
        console.error('Failed to fetch data points data:', err);
        setError('Failed to load data points data');
        setDataPointsData([]);
      } finally {
        setLoading(false);
      }
    };

    fetchDataPointsData();
  }, [timeRange]);

  // Show loading state
  if (loading) {
    return (
      <div
        style={{
          width: '100%',
          height: '400px',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
        }}
      >
        <div style={{ color: '#9CA3AF' }}>Loading data points...</div>
      </div>
    );
  }

  // Show error state
  if (error) {
    return (
      <div
        style={{
          width: '100%',
          height: '400px',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
        }}
      >
        <div style={{ color: '#EF4444' }}>{error}</div>
      </div>
    );
  }

  return (
    <TimeSeriesChart
      width={widthProp}
      height={heightProp}
      data={dataPointsData}
      xAccessor={getDate}
      yAccessor={getDataPoints}
      xFormatter={formatDate}
      colors={greenColors}
      showTooltip={showTooltip}
      hideTooltip={hideTooltip}
      tooltipData={tooltipData}
      tooltipTop={tooltipTop}
      tooltipLeft={tooltipLeft}
    />
  );
};

const DataPointsChartWithTooltip = withTooltip(DataPointsChart);

const DataPoints = () => {
  const [selectedRange, setSelectedRange] = useState('1w');

  return (
    <div className="rounded-lg p-8 min-h-[350px] flex flex-col">
      <div className="flex items-center justify-between gap-3 mb-4">
        <h4 className="text-lg font-semibold text-white">Data Points</h4>
        <TimeRangeToggle
          value={selectedRange}
          onChange={(range) => {
            if (!range || !range.length) return;
            console.log('range', range);
            setSelectedRange(range[0]);
          }}
        />
      </div>
      <div className="flex-1">
        <DataPointsChartWithTooltip timeRange={selectedRange} />
      </div>
    </div>
  );
};

export default DataPoints;
