import * as React from 'react'
import { IconCircleHalf, IconMoonStarsFill, IconSunFill } from '../icons'
import { useState, useEffect } from 'react'
import { Theme, ThemeID } from './theme'

const themes: Theme[] = [
	{
		id: 'light',
		icon: <IconSunFill />,
		label: 'Светлая тэма',
	},
	{
		id: 'dark',
		icon: <IconMoonStarsFill />,
		label: 'Цёмная тэма',
	},
	{
		id: 'auto',
		icon: <IconCircleHalf />,
		label: 'Сістэмная тэма',
	}
]

function ThemeSelector() {
	const [theme, setTheme] = useState<ThemeID | null>(null)

	useEffect(() => {
		setTheme(window.getStoredTheme())
	}, [])

	useEffect(() => {
		if (!theme) {
			return
		}
		window.updateStoredTheme(theme)
	}, [theme])

	useEffect(() => {
		window.addEventListener('storage', e => {
			if (e.key == 'theme') {
				setTheme(window.getStoredTheme())
			}
		})
	}, [])

	useEffect(() => {
		import('bootstrap')
	}, [])

	if (!theme) {
		return <></>
	}

	return <span className='btn-group dropup'>
		<button className="btn btn-sm btn-link text-secondary dropdown-toggle" type="button" data-bs-toggle="dropdown" aria-expanded="false">
			{themes.find(t => t.id == theme).icon}
		</button>
		<ul className="dropdown-menu">
			{themes.map(t => (
				<li><a className={`btn btn-sm dropdown-item ${t.id == theme ? 'active' : ''}`}
					onClick={() => setTheme(t.id)}>{t.icon} {t.label}</a></li>
			))}
		</ul>
	</span>
}

export default ThemeSelector
