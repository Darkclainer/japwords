import { useEffect, useState } from 'react';
import { NavLink, Outlet, useLocation } from 'react-router-dom';

import HealthStatusIcon from '../components/HealthStatusIcon';

export default function Root() {
  const activeBold = ({ isActive }: { isActive: boolean }) => (isActive ? 'font-bold' : '');

  // keep last query to ease navigation between settings and search
  const [lastSearch, setLastSearch] = useState<string>('/search');
  const location = useLocation();
  useEffect(() => {
    if (location.pathname.startsWith('/search')) {
      setLastSearch(location.pathname);
    }
  }, [location]);

  return (
    <>
      <header>
        <nav>
          <div className="flex flex-row gap-x-4 pt-11 pb-5 pl-2 justify-between">
            <div className="flex flex-row gap-x-4 items-center text-2xl text-blue">
              <NavLink
                to={lastSearch == location.pathname ? 'search' : lastSearch}
                className={activeBold}
              >
                Search
              </NavLink>
              <NavLink to="anki" className={activeBold}>
                Anki Settings
              </NavLink>
            </div>
            <HealthStatusIcon className="w-10 h-10" />
          </div>
        </nav>
      </header>
      <main className="flex flex-col flex-1">
        <Outlet></Outlet>
      </main>
    </>
  );
}
