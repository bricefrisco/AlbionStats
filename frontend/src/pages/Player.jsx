import { useParams } from 'react-router-dom';
import Page from '../components/Page';
import PlayerDetails from '../components/PlayerDetails';
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
  return (
    <RegionProvider>
      <Page
        title={`Player ${decodedName ?? 'Profile'}`}
        description={`${detailMessage}.`}
      >
        <div className="mt-2 space-y-3">
            <PlayerDetails region={region} decodedName={decodedName} />
          </div>
      </Page>
    </RegionProvider>
  );
};

export default Player;

