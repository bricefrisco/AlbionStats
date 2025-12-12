import Header from './Header';

const Page = ({ children, title }) => {
  return (
    <div className="bg-black text-white min-h-screen flex flex-col">
      <Header />
      <div className="max-w-7xl mx-auto border-l border-r border-white/15 flex-1 w-full p-4">
        <h1 className="text-3xl font-bold">{title}</h1>
        <div className="mt-200">{children}</div>
      </div>
    </div>
  );
};

export default Page;
