import * as React from 'react';
import TimeSeriesChart from './TimeSeriesChart';
import { formatCompactNumber, formatDate } from '../utils/formatters';

const formatCompactNumberPrecise = (value) => {
  if (value == null) return '0';
  const absValue = Math.abs(value);

  const formatWithSuffix = (val, divisor, suffix) =>
    `${(val / divisor).toFixed(2).replace(/\.?0+$/, '')}${suffix}`;

  if (absValue >= 1_000_000_000) return formatWithSuffix(value, 1_000_000_000, 'b');
  if (absValue >= 1_000_000) return formatWithSuffix(value, 1_000_000, 'm');
  if (absValue >= 1_000) return formatWithSuffix(value, 1_000, 'k');

  return Number(value).toLocaleString(undefined, {
    minimumFractionDigits: 0,
    maximumFractionDigits: 2,
  });
};

const PlayerPvEChart = ({ region, playerName }) => {
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
      )}/pve`,
      {
        signal: controller.signal,
      }
    )
      .then((response) => {
        if (!response.ok) {
          if (response.status === 404) {
            throw new Error('Player not found');
          }
          throw new Error('Unable to load PvE history');
        }
        return response.json();
      })
      .then((payload) => {
        const { Timestamps, PveTotal, PveRoyal, PveOutlands, PveAvalon, PveHellgate, PveCorrupted, PveMists } = payload;
        const normalized = Timestamps.map((timestamp, index) => ({
          timestamp: new Date(timestamp),
          pveTotal: PveTotal[index] ?? 0,
          pveRoyal: PveRoyal[index] ?? 0,
          pveOutlands: PveOutlands[index] ?? 0,
          pveAvalon: PveAvalon[index] ?? 0,
          pveHellgate: PveHellgate[index] ?? 0,
          pveCorrupted: PveCorrupted[index] ?? 0,
          pveMists: PveMists[index] ?? 0,
        }));
        setHistory(normalized);
      })
      .catch((err) => {
        if (err.name === 'AbortError') {
          return;
        }
        console.error(err);
        setError(err.message || 'PvE data could not be loaded');
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
        Select a player to view PvE history.
      </p>
    );
  }

  if (loading) {
    return (
      <p className="text-sm text-gray-400">Loading PvE history…</p>
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
        No PvE history available yet.
      </p>
    );
  }

  const latest = chartContent[chartContent.length - 1];

  return (
    <div className="space-y-4 w-full">
      {/* Total PvE Fame - Full Width */}
      <div className="space-y-2">
        <div className="flex items-center justify-between text-xs uppercase text-gray-400 tracking-wide">
          <span>Total PvE Fame</span>
          {latest && (
            <span className="text-emerald-300">
              {formatCompactNumber(latest.pveTotal)}
            </span>
          )}
        </div>
        <TimeSeriesChart
          data={chartContent}
          xAccessor={(d) => d.timestamp}
          yAccessor={(d) => d.pveTotal}
          xFormatter={formatDate}
          yFormatter={(value) => formatCompactNumberPrecise(value)}
          colors={{
            accentColor: '#10b981',
            accentColorDark: '#059669',
            background: '#1a1a1a',
            background2: '#0a0a0a',
            gridColor: '#374151',
            textColor: '#9CA3AF',
          }}
        />
      </div>

      {/* Remaining categories in 2-column grid */}
      <div className="grid gap-4 md:grid-cols-2">
        <div className="space-y-2">
          <div className="flex items-center justify-between text-xs uppercase text-gray-400 tracking-wide">
            <span>Royal Continent</span>
            {latest && (
              <span className="text-purple-300">
                {formatCompactNumber(latest.pveRoyal)}
              </span>
            )}
          </div>
          <TimeSeriesChart
            data={chartContent}
            xAccessor={(d) => d.timestamp}
            yAccessor={(d) => d.pveRoyal}
            xFormatter={formatDate}
            yFormatter={(value) => formatCompactNumberPrecise(value)}
            colors={{
              accentColor: '#8b5cf6',
              accentColorDark: '#7c3aed',
              background: '#1a1a1a',
              background2: '#0a0a0a',
              gridColor: '#374151',
              textColor: '#9CA3AF',
            }}
          />
        </div>

        <div className="space-y-2">
          <div className="flex items-center justify-between text-xs uppercase text-gray-400 tracking-wide">
            <span>Outlands</span>
            {latest && (
              <span className="text-blue-300">
                {formatCompactNumber(latest.pveOutlands)}
              </span>
            )}
          </div>
          <TimeSeriesChart
            data={chartContent}
            xAccessor={(d) => d.timestamp}
            yAccessor={(d) => d.pveOutlands}
            xFormatter={formatDate}
            yFormatter={(value) => formatCompactNumberPrecise(value)}
            colors={{
              accentColor: '#3b82f6',
              accentColorDark: '#2563eb',
              background: '#1a1a1a',
              background2: '#0a0a0a',
              gridColor: '#374151',
              textColor: '#9CA3AF',
            }}
          />
        </div>

        <div className="space-y-2">
          <div className="flex items-center justify-between text-xs uppercase text-gray-400 tracking-wide">
            <span>Avalon</span>
            {latest && (
              <span className="text-yellow-300">
                {formatCompactNumber(latest.pveAvalon)}
              </span>
            )}
          </div>
          <TimeSeriesChart
            data={chartContent}
            xAccessor={(d) => d.timestamp}
            yAccessor={(d) => d.pveAvalon}
            xFormatter={formatDate}
            yFormatter={(value) => formatCompactNumberPrecise(value)}
            colors={{
              accentColor: '#eab308',
              accentColorDark: '#ca8a04',
              background: '#1a1a1a',
              background2: '#0a0a0a',
              gridColor: '#374151',
              textColor: '#9CA3AF',
            }}
          />
        </div>

        <div className="space-y-2">
          <div className="flex items-center justify-between text-xs uppercase text-gray-400 tracking-wide">
            <span>Hellgate</span>
            {latest && (
              <span className="text-red-300">
                {formatCompactNumber(latest.pveHellgate)}
              </span>
            )}
          </div>
          <TimeSeriesChart
            data={chartContent}
            xAccessor={(d) => d.timestamp}
            yAccessor={(d) => d.pveHellgate}
            xFormatter={formatDate}
            yFormatter={(value) => formatCompactNumberPrecise(value)}
            colors={{
              accentColor: '#ef4444',
              accentColorDark: '#dc2626',
              background: '#1a1a1a',
              background2: '#0a0a0a',
              gridColor: '#374151',
              textColor: '#9CA3AF',
            }}
          />
        </div>

        <div className="space-y-2">
          <div className="flex items-center justify-between text-xs uppercase text-gray-400 tracking-wide">
            <span>Corrupted Dungeons</span>
            {latest && (
              <span className="text-orange-300">
                {formatCompactNumber(latest.pveCorrupted)}
              </span>
            )}
          </div>
          <TimeSeriesChart
            data={chartContent}
            xAccessor={(d) => d.timestamp}
            yAccessor={(d) => d.pveCorrupted}
            xFormatter={formatDate}
            yFormatter={(value) => formatCompactNumberPrecise(value)}
            colors={{
              accentColor: '#f97316',
              accentColorDark: '#ea580c',
              background: '#1a1a1a',
              background2: '#0a0a0a',
              gridColor: '#374151',
              textColor: '#9CA3AF',
            }}
          />
        </div>

        <div className="space-y-2">
          <div className="flex items-center justify-between text-xs uppercase text-gray-400 tracking-wide">
            <span>Mists</span>
            {latest && (
              <span className="text-cyan-300">
                {formatCompactNumber(latest.pveMists)}
              </span>
            )}
          </div>
          <TimeSeriesChart
            data={chartContent}
            xAccessor={(d) => d.timestamp}
            yAccessor={(d) => d.pveMists}
            xFormatter={formatDate}
            yFormatter={(value) => formatCompactNumberPrecise(value)}
            colors={{
              accentColor: '#06b6d4',
              accentColorDark: '#0891b2',
              background: '#1a1a1a',
              background2: '#0a0a0a',
              gridColor: '#374151',
              textColor: '#9CA3AF',
            }}
          />
        </div>
      </div>
    </div>
  );
};

export default PlayerPvEChart;
