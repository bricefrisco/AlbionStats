import Page from './components/Page';
import PlayersTracked from './components/PlayersTracked';
import DataPoints from './components/DataPoints';

const App = () => {
  return (
    <Page
      title="Albion Online Player Statistics"
      description="Search Albion Online players and view their PvP, PvE, and gathering progression data."
    >
      <div className="grid md:grid-cols-2 gap-6 mb-12">
        <PlayersTracked />

        <DataPoints />
      </div>

      <div className="grid md:grid-cols-2 gap-0">
        <div className="p-6 border-r border-b border-white/10">
          <h3 className="text-xl font-semibold mb-3 text-white">PvP Fame</h3>
          <p className="text-gray-300">
            Track player versus player combat performance and kill/death ratios.
          </p>
        </div>

        <div className="p-6 border-b border-white/10">
          <h3 className="text-xl font-semibold mb-3 text-white">PvE Fame</h3>
          <p className="text-gray-300">
            Monitor player versus environment progression and achievements.
          </p>
        </div>

        <div className="p-6 border-r border-white/10">
          <h3 className="text-xl font-semibold mb-3 text-white">
            Gathering Fame
          </h3>
          <p className="text-gray-300">
            View resource collection statistics across fiber, hide, ore, rock,
            and wood.
          </p>
        </div>

        <div className="p-6">
          <h3 className="text-xl font-semibold mb-3 text-white">
            Guild History
          </h3>
          <p className="text-gray-300">
            Explore guild memberships and alliance affiliations over time.
          </p>
        </div>
      </div>
    </Page>
  );
};

export default App;
