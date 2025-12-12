import { useState } from 'react';

import Page from './components/Page';
import PlayersTracked from './components/PlayersTracked';
import TimeRangeToggle from './components/TimeRangeToggle';

const App = () => {
  const [selectedRange, setSelectedRange] = useState('1w');

  return (
    <Page
      title="Albion Online Player Statistics"
      description="Search Albion Online players and view their statistics, guild affiliations, and PvE progression data."
    >
      <div className="grid md:grid-cols-2 gap-6 mb-12">
        <div className="rounded-lg p-8 border border-white/10 min-h-[350px] flex flex-col">
          <div className="flex items-center justify-between gap-3 mb-4">
            <h4 className="text-lg font-semibold text-white">
              Players Tracked
            </h4>
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
            <PlayersTracked timeRange={selectedRange} />
          </div>
        </div>

        <div className="rounded-lg p-8 border border-white/10 min-h-[350px] flex flex-col">
          <h4 className="text-lg font-semibold text-white mb-4">Data Points</h4>
          <div className="flex items-end justify-between flex-1 space-x-2">
            <div
              className="bg-green-400/20 rounded-sm flex-1"
              style={{ height: '30%' }}
            ></div>
            <div
              className="bg-green-400/20 rounded-sm flex-1"
              style={{ height: '50%' }}
            ></div>
            <div
              className="bg-green-400/20 rounded-sm flex-1"
              style={{ height: '75%' }}
            ></div>
            <div
              className="bg-green-400/20 rounded-sm flex-1"
              style={{ height: '85%' }}
            ></div>
            <div
              className="bg-green-400/20 rounded-sm flex-1"
              style={{ height: '95%' }}
            ></div>
            <div
              className="bg-green-400/20 rounded-sm flex-1"
              style={{ height: '100%' }}
            ></div>
          </div>
        </div>
      </div>

      <div className="grid md:grid-cols-2 gap-6">
        <div className="rounded-lg p-6 border border-white/10">
          <h3 className="text-xl font-semibold mb-3 text-white">PvP Fame</h3>
          <p className="text-gray-300">
            Track player versus player combat performance and kill/death ratios.
          </p>
        </div>

        <div className="rounded-lg p-6 border border-white/10">
          <h3 className="text-xl font-semibold mb-3 text-white">PvE Fame</h3>
          <p className="text-gray-300">
            Monitor player versus environment progression and achievements.
          </p>
        </div>

        <div className="rounded-lg p-6 border border-white/10">
          <h3 className="text-xl font-semibold mb-3 text-white">
            Gathering Fame
          </h3>
          <p className="text-gray-300">
            View resource collection statistics across fiber, hide, ore, rock,
            and wood.
          </p>
        </div>

        <div className="rounded-lg p-6 border border-white/10">
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
