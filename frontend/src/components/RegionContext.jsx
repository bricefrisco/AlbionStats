import * as React from 'react';

const RegionContext = React.createContext();

export const RegionProvider = ({ children }) => {
  const [region, setRegion] = React.useState('americas'); // Default to Americas

  const value = {
    region,
    setRegion,
  };

  return (
    <RegionContext.Provider value={value}>{children}</RegionContext.Provider>
  );
};

export const useRegion = () => {
  const context = React.useContext(RegionContext);
  if (!context) {
    throw new Error('useRegion must be used within a RegionProvider');
  }
  return context;
};
