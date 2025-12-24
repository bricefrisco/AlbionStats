import * as React from 'react';
import { useParams } from 'react-router-dom';
import Page from '../components/Page';
import { RegionProvider } from '../components/RegionContext';

const Player = () => {
  const { region, name } = useParams();
  const decodedName = name ? decodeURIComponent(name) : undefined;
  const trimmedName = decodedName?.trim() || 'Unknown';
  const capitalize = (value) => {
    if (!value) return '';
    return `${value.charAt(0).toUpperCase()}${value.slice(1)}`;
  };
  const regionLabel = region ? capitalize(region) : 'Unknown';
  const detailMessage = `Detailed view for ${trimmedName} on the ${regionLabel} server`;
  const [playerStats, setPlayerStats] = React.useState(null);
  const [isLoadingStats, setIsLoadingStats] = React.useState(false);
  const [statsError, setStatsError] = React.useState(null);

  const formatNumber = (value) => {
    if (value == null) return '0';
    return Intl.NumberFormat('en-US').format(value);
  };

  React.useEffect(() => {
    if (!region || !decodedName) {
      setPlayerStats(null);
      setStatsError(null);
      return undefined;
    }

    const controller = new AbortController();
    setIsLoadingStats(true);
    setStatsError(null);

    const playerUrl = `https://api.bricefrisco.com/albionstats/v1/players/${region}/${encodeURIComponent(
      decodedName
    )}`;

    fetch(playerUrl, {
      signal: controller.signal,
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error('Player stats could not be loaded');
        }
        return response.json();
      })
      .then((data) => {
        setPlayerStats(data);
      })
      .catch((error) => {
        if (error.name === 'AbortError') {
          return;
        }
        console.error(error);
        setStatsError('Unable to load player stats right now');
        setPlayerStats(null);
      })
      .finally(() => {
        setIsLoadingStats(false);
      });

    return () => {
      controller.abort();
    };
  }, [decodedName, region]);

  return (
    <RegionProvider>
      <Page
        title={`Player ${decodedName ?? 'Profile'}`}
        description={`${detailMessage}.`}
      >
        <div className="rounded-lg border border-white/10 bg-white/5 px-6 pt-2 pb-4">
          <div className="mt-2 space-y-3">
            {isLoadingStats && (
              <p className="text-sm text-gray-400">Loading player stats…</p>
            )}
            {statsError && (
              <p className="text-sm text-red-400">{statsError}</p>
            )}
            {playerStats && !statsError && (
              <div className="flex flex-col gap-6 md:flex-row">
                <div className="flex-1 space-y-3">
                  <div>
                    <p className="text-xs uppercase text-gray-400">Kill Fame</p>
                    <p className="text-white">
                      {formatNumber(playerStats.KillFame)}
                    </p>
                  </div>
                  <div>
                    <p className="text-xs uppercase text-gray-400">Death Fame</p>
                    <p className="text-white">
                      {formatNumber(playerStats.DeathFame)}
                    </p>
                  </div>
                  <div>
                    <p className="text-xs uppercase text-gray-400">Fame Ratio</p>
                    <p className="text-white">
                      {playerStats.FameRatio != null
                        ? Number(playerStats.FameRatio).toFixed(2)
                        : 'N/A'}
                    </p>
                  </div>
                </div>
                {(playerStats.GuildName || playerStats.AllianceName) && (
                  <div className="flex-1 space-y-3">
                    {playerStats.GuildName && (
                      <div>
                        <p className="text-xs uppercase text-gray-400">Guild</p>
                        <p className="text-white">{playerStats.GuildName}</p>
                      </div>
                    )}
                    {playerStats.AllianceName && (
                      <div>
                        <p className="text-xs uppercase text-gray-400">Alliance</p>
                        <p className="text-white">{playerStats.AllianceName}</p>
                      </div>
                    )}
                  </div>
                )}
              </div>
            )}
          </div>
        </div>
      </Page>
    </RegionProvider>
  );
};

export default Player;

