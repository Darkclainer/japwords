import { NavLink, Outlet } from 'react-router-dom';

export default function Anki() {
  return (
    <>
      <p>ANKI ROOT</p>
      <NavLink to="config">config</NavLink>
      <NavLink to="state">state</NavLink>
      <Outlet></Outlet>
    </>
  );
}
