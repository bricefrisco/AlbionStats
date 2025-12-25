import * as React from 'react';
import TimeSeriesChart from './TimeSeriesChart';
import { formatCompactNumber, formatDate } from '../utils/formatters';

const PlayerPvPChart = ({ region, playerName }) => {
  const [history, setHistory] = React.useState([]);
  const [loading, setLoading] = React.useState(false);
  const [error, setError] = React.useState(null);

  React.useEffect(() => {
    if (!region || !playerName) {
      setHistory([]);
      return;
    }

    const controller = new AbortController();
    setLoading(true);
    setError(null);

    fetch(
      `https://api.bricefrisco.com/albionstats/v1/players/${region}/${encodeURIComponent(
        playerName
      )}/pvp`,
      {
        signal: controller.signal,
      }
    )
      .then((response) => {
        if (!response.ok) {
          if (response.status === 404) {
            throw new Error('Player not found');
          }
          throw new Error('Unable to load PvP history');
        }
        return response.json();
      })
      .then((payload) => {
        const { Timestamps, KillFame, DeathFame, FameRatio } = payload;
        const normalized = Timestamps.map((timestamp, index) => ({
          timestamp: new Date(timestamp),
          killFame: KillFame[index] ?? 0,
          deathFame: DeathFame[index] ?? 0,
          fameRatio: FameRatio[index] ?? null,
        }));
        setHistory(normalized);
      })
      .catch((err) => {
        if (err.name === 'AbortError') {
          return;
        }
        console.error(err);
        setError(err.message || 'PvP data could not be loaded');
        setHistory([]);
      })
      .finally(() => {
        setLoading(false);
      });

    return () => {
      controller.abort();
    };
  }, [region, playerName]);

  const chartContent = React.useMemo(() => {
    if (!history || history.length === 0) return [];
    return history;
  }, [history]);

  if (!region || !playerName) {
    return (
      <p className="text-sm text-gray-400">
        Select a player to view PvP history.
      </p>
    );
  }

  if (loading) {
    return (
      <p className="text-sm text-gray-400">Loading PvP history…</p>
    );
  }

  if (error) {
    return (
      <p className="text-sm text-red-400">
        {error}
      </p>
    );
  }

  if (!chartContent.length) {
    return (
      <p className="text-sm text-gray-400">
        No PvP history available yet.
      </p>
    );
  }

  const latest = chartContent[chartContent.length - 1];

  return (
    <div className="space-y-4 w-full">
      <div className="flex items-center justify-between text-xs uppercase text-gray-400 tracking-wide">
        <span>PvP Kill Fame</span>
        {latest && (
            <span className="text-white">
              Latest: {formatCompactNumber(latest.killFame)}
            </span>
        )}
      </div>
      <TimeSeriesChart
        data={chartContent}
        xAccessor={(d) => d.timestamp}
        yAccessor={(d) => d.killFame}
        xFormatter={formatDate}
        yFormatter={(value) => formatCompactNumber(value)}
      />
      <div className="flex items-center justify-between text-xs uppercase text-gray-400 tracking-wide">
        <span>PvP Death Fame</span>
        {latest && (
            <span className="text-white">
              Latest: {formatCompactNumber(latest.deathFame)}
            </span>
        )}
      </div>
      <TimeSeriesChart
        data={chartContent}
        xAccessor={(d) => d.timestamp}
        yAccessor={(d) => d.deathFame}
        xFormatter={formatDate}
        yFormatter={(value) => formatCompactNumber(value)}
      />
    </div>
  );
};

export default PlayerPvPChart;

