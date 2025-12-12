import { Input } from '@base-ui/react/input';

const Header = () => {
  return (
    <header className="border-b border-white/15">
      <div className="max-w-7xl mx-auto border-l border-r border-white/15 flex items-center">
        <div className="py-3 px-5 text-lg font-bold">AlbionStats</div>
        <div className="flex-1 flex justify-center">
          <div className="relative">
            <Input
              placeholder="Search players..."
              className="bg-zinc-900 border border-white/15 rounded-full pl-10 pr-5 py-1 min-w-96 focus:border-blue-400 focus:outline-none placeholder-gray-400"
            />
            <img
              src="/search.svg"
              alt="Search"
              className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 brightness-0 invert"
            />
          </div>
        </div>
        <a
          href="https://github.com/bricefrisco/AlbionStats"
          target="_blank"
          rel="noopener noreferrer"
          className="ml-auto px-5 flex items-center gap-2 hover:underline"
        >
          <img
            src="/github.svg"
            alt="GitHub"
            className="w-5 h-5 brightness-0 invert"
          />
          <span className="text-md">Github</span>
        </a>
      </div>
    </header>
  );
};

export default Header;
