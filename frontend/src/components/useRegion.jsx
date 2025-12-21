import * as React from 'react';
import { RegionContext } from './region-context';

export const useRegion = () => {
  const context = React.useContext(RegionContext);
  if (!context) {
    throw new Error('useRegion must be used within a RegionProvider');
  }
  return context;
};

export default useRegion;

