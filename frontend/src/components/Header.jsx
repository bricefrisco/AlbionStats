import { Link } from 'react-router-dom';
import Search from './Search';
import Select from './Select';
import { useRegion } from './useRegion';

const regions = [
  { label: 'Americas', value: 'americas' },
  { label: 'Europe', value: 'europe' },
  { label: 'Asia', value: 'asia' },
];

const Header = () => {
  const { region, setRegion } = useRegion();

  return (
    <header className="sticky top-0 z-10 border-b border-white/15 bg-black/80 backdrop-blur-sm">
      <div className="max-w-7xl mx-auto border-l border-r border-white/15 grid grid-cols-3 items-center">
        <Link
          to="/"
          className="py-3 px-5 text-lg font-bold hover:underline"
        >
          AlbionStats
        </Link>
        <div className="flex justify-center">
          <Search />
        </div>
        <div className="flex items-center justify-end gap-4 px-5">
          <Select items={regions} value={region} onValueChange={setRegion} />
          <a
            href="https://github.com/bricefrisco/AlbionStats"
            target="_blank"
            rel="noopener noreferrer"
            className="flex items-center gap-2 hover:underline"
          >
            <img
              src="/github.svg"
              alt="GitHub"
              className="w-5 h-5 brightness-0 invert"
            />
            <span className="text-md">Github</span>
          </a>
        </div>
      </div>
    </header>
  );
};

export default Header;
