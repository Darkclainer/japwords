import { clsx } from 'clsx';
import { NavLink, Outlet } from 'react-router-dom';

function Tab(props: { title: string; to: string; last?: boolean; active?: boolean }) {
  const { title, to, last } = props;
  return (
    <NavLink
      to={to}
      className={({ isActive }) =>
        clsx(
          'pt-9 pb-5 px-6',
          'text-blue',
          isActive ? 'font-bold bg-gray' : 'bg-mid-gray shadow-tab-inner',
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
      <div className="flex flex-col shadow-md rounded-md bg-mid-gray">
        <div className="flex flex-row text-2xl text-slate-500">
          <Tab to="connection-settings" title="Anki Connect" />
          <Tab to="user-settings" title="Anki User Settings" last />
        </div>
        <div className="bg-gray pt-8 pb-8 px-10">
          <Outlet></Outlet>
        </div>
      </div>
    </>
  );
}
