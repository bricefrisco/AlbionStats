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

  return (
    <RegionProvider>
      <Page
        title={`Player ${decodedName ?? 'Profile'}`}
        description={`${detailMessage}.`}
      >
        <div className="rounded-lg border border-white/10 bg-white/5 p-6">
          <p className="text-gray-300">
            This is a placeholder for{' '}
            <span className="text-white">
              {decodedName ?? 'the selected player'}
            </span>{' '}
            on <span className="text-white">{region ?? 'the selected region'}</span>.
          </p>
          <p className="mt-4 text-lg text-gray-300">{detailMessage}</p>
          <p className="mt-3 text-sm text-gray-400">
            Plug in the actual player stats and navigation controls here.
          </p>
        </div>
      </Page>
    </RegionProvider>
  );
};

export default Player;

