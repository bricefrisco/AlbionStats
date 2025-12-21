import React, { useEffect, useState } from 'react';
import { withTooltip } from '@visx/tooltip';
import TimeSeriesChart from './TimeSeriesChart';
import TimeRangeToggle from './TimeRangeToggle';

// API endpoint configuration
const API_BASE_URL = 'https://api.bricefrisco.com';
const METRICS_ENDPOINT = '/albionstats/v1/metrics/players_total';

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

// PlayersTrackedChart component for Fast Refresh compatibility
const PlayersTrackedChart = ({
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
  const [playerData, setPlayerData] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Fetch player data from API
  useEffect(() => {
    const fetchPlayerData = async () => {
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
          players: item.value,
        }));

        setPlayerData(transformedData);
        setError(null);
      } catch (err) {
        console.error('Failed to fetch player data:', err);
        setError('Failed to load player data');
        setPlayerData([]);
      } finally {
        setLoading(false);
      }
    };

    fetchPlayerData();
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
        <div style={{ color: '#9CA3AF' }}>Loading player data...</div>
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
      data={playerData}
      xAccessor={getDate}
      yAccessor={getPlayerCount}
      xFormatter={formatDate}
      showTooltip={showTooltip}
      hideTooltip={hideTooltip}
      tooltipData={tooltipData}
      tooltipTop={tooltipTop}
      tooltipLeft={tooltipLeft}
    />
  );
};

const PlayersTrackedChartWithTooltip = withTooltip(PlayersTrackedChart);

const PlayersTracked = () => {
  const [selectedRange, setSelectedRange] = useState('1w');

  return (
    <div className="rounded-lg p-8 min-h-[350px] flex flex-col">
      <div className="flex items-center justify-between gap-3 mb-4">
        <h4 className="text-lg font-semibold text-white">Players Tracked</h4>
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
        <PlayersTrackedChartWithTooltip timeRange={selectedRange} />
      </div>
    </div>
  );
};

export default PlayersTracked;
