import * as React from 'react';
import { Autocomplete } from '@base-ui/react/autocomplete';
import { useRegion } from './RegionContext';

export default function Search() {
  const { region } = useRegion();
  const [searchValue, setSearchValue] = React.useState('');
  const [isLoading, setIsLoading] = React.useState(false);
  const [searchResults, setSearchResults] = React.useState([]);
  const [error, setError] = React.useState(null);

  const formatGuildAlliance = (guildName, allianceName) => {
    if (allianceName && guildName) {
      return `[${allianceName}] ${guildName}`;
    } else if (guildName) {
      return guildName;
    }
    return '';
  };

  React.useEffect(() => {
    if (!searchValue.trim() || searchValue.trim().length < 3) {
      setSearchResults([]);
      setIsLoading(false);
      return undefined;
    }

    setIsLoading(true);
    setError(null);

    async function fetchPlayers() {
      try {
        const response = await fetch(
          `https://api.bricefrisco.com/albionstats/v1/search/${region}/${encodeURIComponent(searchValue)}`
        );
        if (!response.ok) {
          throw new Error('Failed to fetch players');
        }
        const data = await response.json();
        setSearchResults(data.players || []);
      } catch (err) {
        console.error(err);
        setError('Failed to search players. Please try again.');
        setSearchResults([]);
      } finally {
        setIsLoading(false);
      }
    }

    const timeoutId = setTimeout(fetchPlayers, 300);

    return () => {
      clearTimeout(timeoutId);
    };
  }, [searchValue, region]);

  let status = `${searchResults.length} result${searchResults.length === 1 ? '' : 's'} found`;
  if (isLoading) {
    status = 'Searching...';
  } else if (error) {
    status = error;
  } else if (searchResults.length === 0 && searchValue.trim()) {
    status = `No players found for "${searchValue}"`;
  }

  const shouldRenderPopup = searchValue.trim().length >= 3;

  return (
    <Autocomplete.Root
      items={searchResults}
      value={searchValue}
      onValueChange={setSearchValue}
      itemToStringValue={(item) => item.name}
      filter={null}
    >
      <div className="relative">
        <Autocomplete.Input
          placeholder="Search players..."
          className="bg-zinc-900 border border-white/15 rounded-full pl-10 pr-5 py-1 min-w-96 focus:border-blue-400 focus:outline-none placeholder-gray-400"
        />
        <img
          src="/search.svg"
          alt="Search"
          className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 brightness-0 invert"
        />
      </div>

      {shouldRenderPopup && (
        <Autocomplete.Portal>
          <Autocomplete.Positioner
            className="absolute z-50 mt-1"
            sideOffset={4}
            align="start"
          >
            <Autocomplete.Popup className="bg-zinc-900 border border-white/15 rounded-lg shadow-lg w-full min-w-96">
              <Autocomplete.Status className="px-3 py-2 text-sm text-gray-400">
                {status}
              </Autocomplete.Status>
              <Autocomplete.List className="max-h-64 overflow-auto">
                {(player) => (
                  <Autocomplete.Item
                    key={player.player_id}
                    className="px-3 py-2 hover:bg-white/10 cursor-pointer text-white"
                    value={player}
                  >
                    <div className="flex flex-col">
                      <div className="font-medium">{player.name}</div>
                      {formatGuildAlliance(
                        player.guild_name,
                        player.alliance_name
                      ) && (
                        <div className="text-xs text-gray-400">
                          {formatGuildAlliance(
                            player.guild_name,
                            player.alliance_name
                          )}
                        </div>
                      )}
                    </div>
                  </Autocomplete.Item>
                )}
              </Autocomplete.List>
            </Autocomplete.Popup>
          </Autocomplete.Positioner>
        </Autocomplete.Portal>
      )}
    </Autocomplete.Root>
  );
}
