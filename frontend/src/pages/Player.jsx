import * as React from 'react';
import { useParams } from 'react-router-dom';
import Page from '../components/Page';
import PlayerDetail from '../components/PlayerDetail';
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
        <div className="mt-2 space-y-3">
            {isLoadingStats && (
              <p className="text-sm text-gray-400">Loading player stats…</p>
            )}
            {statsError && (
              <p className="text-sm text-red-400">{statsError}</p>
            )}
            {playerStats && !statsError && (
              <div className="columns-1 md:columns-2">
                {(playerStats.GuildName || playerStats.AllianceName) && (
                  <PlayerDetail title="Guild">
                    <div className="grid grid-cols-2 gap-3 items-start">
                      {playerStats.GuildName && (
                        <PlayerDetail.Item label="Guild">
                          {playerStats.GuildName}
                        </PlayerDetail.Item>
                      )}
                      {playerStats.AllianceName && (
                        <PlayerDetail.Item label="Alliance">
                          {playerStats.AllianceName}
                        </PlayerDetail.Item>
                      )}
                    </div>
                  </PlayerDetail>
                )}
                <PlayerDetail title="Crafting">
                  <PlayerDetail.Item label="Crafting Fame">
                    {formatNumber(playerStats.CraftingTotal)}
                  </PlayerDetail.Item>
                </PlayerDetail>
                <PlayerDetail title="Player Versus Environment">
                  <PlayerDetail.Item label="Total">
                    {formatNumber(playerStats.PveTotal)}
                  </PlayerDetail.Item>
                  <div className="grid grid-cols-2 gap-3 items-start">
                    <PlayerDetail.Item label="Royal Continent">
                      {formatNumber(playerStats.PveRoyal)}
                    </PlayerDetail.Item>
                    <PlayerDetail.Item label="Outlands">
                      {formatNumber(playerStats.PveOutlands)}
                    </PlayerDetail.Item>
                  </div>
                  <div className="grid grid-cols-2 gap-3 items-start">
                    <PlayerDetail.Item label="Avalonian Roads">
                      {formatNumber(playerStats.PveAvalon)}
                    </PlayerDetail.Item>
                    <PlayerDetail.Item label="Mists">
                      {formatNumber(playerStats.PveMists)}
                    </PlayerDetail.Item>
                  </div>
                  <div className="grid grid-cols-2 gap-3 items-start">
                    <PlayerDetail.Item label="Hellgates">
                      {formatNumber(playerStats.PveHellgate)}
                    </PlayerDetail.Item>
                    <PlayerDetail.Item label="Corrupted Dungeons">
                      {formatNumber(playerStats.PveCorrupted)}
                    </PlayerDetail.Item>
                  </div>
                </PlayerDetail>
                <PlayerDetail title="Player Versus Player">
                  <div className="grid grid-cols-2 gap-3 items-start">
                    <PlayerDetail.Item label="Kill Fame">
                      {formatNumber(playerStats.KillFame)}
                    </PlayerDetail.Item>
                    <PlayerDetail.Item label="Death Fame">
                      {formatNumber(playerStats.DeathFame)}
                    </PlayerDetail.Item>
                  </div>
                  <PlayerDetail.Item label="Fame Ratio">
                    {playerStats.FameRatio != null
                      ? Number(playerStats.FameRatio).toFixed(2)
                      : 'N/A'}
                  </PlayerDetail.Item>
                </PlayerDetail>
                <PlayerDetail title="Gathering">
                  <PlayerDetail.Item label="Total">
                    {formatNumber(playerStats.GatherAllTotal)}
                  </PlayerDetail.Item>
                  <div className="grid grid-cols-2 gap-3 items-start">
                    <PlayerDetail.Item label="Royal Continent">
                      {formatNumber(playerStats.GatherAllRoyal)}
                    </PlayerDetail.Item>
                    <PlayerDetail.Item label="Outlands">
                      {formatNumber(playerStats.GatherAllOutlands)}
                    </PlayerDetail.Item>
                  </div>
                  <PlayerDetail.Item label="Avalonian Roads">
                    {formatNumber(playerStats.GatherAllAvalon)}
                  </PlayerDetail.Item>
                </PlayerDetail>
              </div>
            )}
          </div>
      </Page>
    </RegionProvider>
  );
};

export default Player;

