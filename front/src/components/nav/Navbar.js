import './style.css'

import { NavLink } from 'react-router-dom'
import BtnDarkMode from '../btnDarkMode/BtnDarkMode'

const Navbar = () => {
	const activeLink = 'nav-list__link nav-list__link--active'
	const normalLink = 'nav-list__link nav-list__link'

	return (
		<nav className='nav'>
			<div className='container'>
				<div className='nav-row'>
					<a
						href='https://github.com/xLeSHka/yandexmicroserver'
						target='_blank'
						rel='noopener noreferrer'
						className='logo'
					>
						<strong>Distributed</strong> calculator
					</a>

					<BtnDarkMode />

					<ul className='nav-list'>
						<li className='nav-list__item'>
							<NavLink
								to='/expressions'
								className={({ isActive }) =>
									isActive ? activeLink : normalLink
								}
							>
								Expressions
							</NavLink>
						</li>
						<li className='nav-list__item'>
							<NavLink
								to='/operations'
								className={({ isActive }) =>
									isActive ? activeLink : normalLink
								}
							>
								Operations
							</NavLink>
						</li>
						<li className='nav-list__item'>
							<NavLink
								to='/agents'
								className={({ isActive }) =>
									isActive ? activeLink : normalLink
								}
							>
								Agents
							</NavLink>
						</li>
					</ul>
				</div>
			</div>
		</nav>
	)
}

export default Navbar
