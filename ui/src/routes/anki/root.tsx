import { clsx } from 'clsx';
import { NavLink, Outlet } from 'react-router-dom';

function Tab(props: { title: string; to: string; last?: boolean; active?: boolean }) {
  const { title, to, last } = props;
  return (
    <NavLink
      to={to}
      className={({ isActive }) =>
        clsx(
          'px-6 pb-5 pt-9 text-blue',
          isActive ? 'bg-gray font-bold' : 'bg-mid-gray shadow-tab-inner',
          last && 'grow',
        )
      }
    >
      {title}
    </NavLink>
  );
}

export default function AnkiRoot() {
  return (
    <>
      <div className="mb-8 flex flex-col rounded-md bg-mid-gray shadow-md">
        <div className="text-slate-500 flex flex-row text-2xl">
          <Tab to="connection-settings" title="Anki Connect" />
          <Tab to="user-settings" title="Anki User Settings" last />
        </div>
        <div className="bg-gray px-10 pb-8 pt-8">
          <Outlet></Outlet>
        </div>
      </div>
    </>
  );
}
