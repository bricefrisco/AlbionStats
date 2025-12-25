import * as React from 'react';
import { Tabs } from '@base-ui/react/tabs';
import PlayerPvPChart from './PlayerPvPChart';

const tabKeys = ['pvp', 'pve', 'gathering', 'crafting'];
const baseTabClass =
  'relative inline-flex items-center justify-center whitespace-nowrap rounded-md px-2 py-1 text-xs font-medium transition-colors focus-visible:z-10 focus-visible:outline focus-visible:outline-2 focus-visible:outline-sky-500 focus-visible:outline-offset-[-2px]';
const activeTabClass =
  'data-[active]:text-white data-[active]:shadow-[inset_0_0_0_1px_rgba(255,255,255,0.15)] data-[active]:font-semibold data-[active]:ring-1 data-[active]:ring-white/25 data-[active]:ring-inset';
const inactiveTabClass = 'text-gray-300 hover:text-white';

const PlayerCharts = ({ region, playerName }) => {
  const [activeTab, setActiveTab] = React.useState('pvp');

  return (
    <Tabs.Root
      className="rounded-lg border border-white/10 bg-transparent shadow-sm"
      value={activeTab}
      onValueChange={setActiveTab}
    >
      <Tabs.List className="relative z-0 flex gap-1 px-2 py-1.5 bg-white/5 shadow-[inset_0_-1px_rgba(255,255,255,0.08)]">
        {tabKeys.map((key) => (
          <Tabs.Tab
            key={key}
            className={`${baseTabClass} ${inactiveTabClass} ${activeTabClass}`}
            value={key}
          >
            {key === 'pvp' && 'PvP'}
            {key === 'pve' && 'PvE'}
            {key === 'gathering' && 'Gathering'}
            {key === 'crafting' && 'Crafting'}
          </Tabs.Tab>
        ))}
        <Tabs.Indicator className="pointer-events-none absolute left-0 top-1/2 z-[-1] h-8 -translate-y-1/2 rounded-md bg-white/10 backdrop-blur transition-[width,transform] duration-200 ease-in-out" />
      </Tabs.List>
      <Tabs.Panel
        className="relative flex min-h-[8rem] items-center justify-center bg-transparent p-4 text-sm text-gray-100 focus-visible:outline focus-visible:outline-2 focus-visible:outline-sky-500 focus-visible:outline-offset-[-2px]"
        value="pvp"
      >
        {activeTab === 'pvp' && (
          <PlayerPvPChart region={region} playerName={playerName} />
        )}
      </Tabs.Panel>
      <Tabs.Panel
        className="relative flex min-h-[8rem] items-center justify-center bg-transparent p-4 text-sm text-gray-100 focus-visible:outline focus-visible:outline-2 focus-visible:outline-sky-500 focus-visible:outline-offset-[-2px]"
        value="pve"
      >
        {activeTab === 'pve' && 'PvE charts will go here.'}
      </Tabs.Panel>
      <Tabs.Panel
        className="relative flex min-h-[8rem] items-center justify-center bg-transparent p-4 text-sm text-gray-100 focus-visible:outline focus-visible:outline-2 focus-visible:outline-sky-500 focus-visible:outline-offset-[-2px]"
        value="gathering"
      >
        {activeTab === 'gathering' && 'Gathering charts will go here.'}
      </Tabs.Panel>
      <Tabs.Panel
        className="relative flex min-h-[8rem] items-center justify-center bg-transparent p-4 text-sm text-gray-100 focus-visible:outline focus-visible:outline-2 focus-visible:outline-sky-500 focus-visible:outline-offset-[-2px]"
        value="crafting"
      >
        {activeTab === 'crafting' && 'Crafting charts will go here.'}
      </Tabs.Panel>
    </Tabs.Root>
  );
};

export default PlayerCharts;

