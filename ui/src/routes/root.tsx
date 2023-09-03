import { useEffect, useState } from 'react';
import { Outlet, NavLink, useLocation } from 'react-router-dom';

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
      <nav className="flex flex-row gap-x-4 text-blue">
        <NavLink
          to={lastSearch == location.pathname ? 'search' : lastSearch}
          className={activeBold}
        >
          Search
        </NavLink>
        <NavLink to="anki" className={activeBold}>
          Anki Settings
        </NavLink>
      </nav>
      <Outlet></Outlet>
    </>
  );
}
