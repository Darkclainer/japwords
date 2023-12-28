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
          <div className="flex flex-row justify-between gap-x-4 pb-5 pl-2 pt-11">
            <div className="flex flex-row items-center gap-x-4 text-2xl text-blue">
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
            <HealthStatusIcon className="h-10 w-10" />
          </div>
        </nav>
      </header>
      <main className="flex flex-1 flex-col">
        <Outlet></Outlet>
      </main>
    </>
  );
}
