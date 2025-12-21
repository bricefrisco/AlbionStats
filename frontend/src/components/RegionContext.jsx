import * as React from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { RegionContext } from './region-context';

const REGION_STORAGE_KEY = 'albion-region';
const VALID_REGIONS = ['americas', 'europe', 'asia'];

const getStoredRegion = () => {
  try {
    if (typeof window === 'undefined') return null;
    return window.localStorage.getItem(REGION_STORAGE_KEY);
  } catch {
    return null;
  }
};

const setStoredRegion = (value) => {
  try {
    if (typeof window === 'undefined') return;
    window.localStorage.setItem(REGION_STORAGE_KEY, value);
  } catch {
    // ignore
  }
};

const parseRegionFromPath = (pathname) => {
  const match = /^\/players\/(americas|europe|asia)(?:\/|$)/.exec(pathname);
  return match ? match[1] : null;
};

export const RegionProvider = ({ children }) => {
  const location = useLocation();
  const navigate = useNavigate();

  const initialRegion =
    parseRegionFromPath(location.pathname) ?? getStoredRegion() ?? 'americas';
  const skipSyncRef = React.useRef(false);
  const [region, setRegionState] = React.useState(initialRegion);

  React.useEffect(() => {
    const pathRegion = parseRegionFromPath(location.pathname);
    if (skipSyncRef.current) {
      if (!pathRegion) {
        skipSyncRef.current = false;
      }
      return;
    }

    if (pathRegion && pathRegion !== region) {
      setRegionState(pathRegion);
      setStoredRegion(pathRegion);
    }
  }, [location.pathname, region]);

  const setRegion = React.useCallback(
    (nextRegion) => {
      if (!VALID_REGIONS.includes(nextRegion) || nextRegion === region) {
        return;
      }

      skipSyncRef.current = true;
      setRegionState(nextRegion);
      setStoredRegion(nextRegion);
      navigate('/');
    },
    [navigate, region]
  );

  const value = {
    region,
    setRegion,
  };

  return (
    <RegionContext.Provider value={value}>{children}</RegionContext.Provider>
  );
};
