import Page from './components/Page';

const App = () => {
  return (
    <Page
      title="AlbionStats - Albion Online Player Statistics"
      description="Search Albion Online players and view their statistics, guild affiliations, and PvE progression data."
    >
      <div className="text-center py-8">
        <div className="max-w-2xl mx-auto mb-12">
          <p className="text-lg text-gray-400">
            Enter a player name in the search bar above to view detailed
            statistics
          </p>
        </div>

        <div className="grid md:grid-cols-2 gap-6 mb-12 max-w-md mx-auto">
          <div className="rounded-lg p-6 border border-white/10">
            <div className="text-3xl font-bold text-blue-400 mb-2">12,847</div>
            <div className="text-sm text-gray-400">Players Tracked</div>
          </div>

          <div className="rounded-lg p-6 border border-white/10">
            <div className="text-3xl font-bold text-green-400 mb-2">2.1M</div>
            <div className="text-sm text-gray-400">Data Points</div>
          </div>
        </div>

        <div className="grid md:grid-cols-2 gap-6">
          <div className="rounded-lg p-6 border border-white/10">
            <h3 className="text-xl font-semibold mb-3 text-blue-400">
              PvP Fame
            </h3>
            <p className="text-gray-300">
              Track player versus player combat performance and kill/death
              ratios.
            </p>
          </div>

          <div className="rounded-lg p-6 border border-white/10">
            <h3 className="text-xl font-semibold mb-3 text-green-400">
              PvE Fame
            </h3>
            <p className="text-gray-300">
              Monitor player versus environment progression and achievements.
            </p>
          </div>

          <div className="rounded-lg p-6 border border-white/10">
            <h3 className="text-xl font-semibold mb-3 text-yellow-400">
              Gathering Fame
            </h3>
            <p className="text-gray-300">
              View resource collection statistics across fiber, hide, ore, rock,
              and wood.
            </p>
          </div>

          <div className="rounded-lg p-6 border border-white/10">
            <h3 className="text-xl font-semibold mb-3 text-purple-400">
              Guild History
            </h3>
            <p className="text-gray-300">
              Explore guild memberships and alliance affiliations over time.
            </p>
          </div>
        </div>
      </div>
    </Page>
  );
};

export default App;
