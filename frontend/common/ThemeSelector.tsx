import * as React from 'react'
import Dropdown from 'react-bootstrap/Dropdown'
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
	const [theme, setTheme] = useState<ThemeID|null>(null)

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

	if (!theme) {
		return <></>
	}

	return (
		<Dropdown as="span" drop='up'>
		  <Dropdown.Toggle size="sm" variant="link" id="theme-selector" className='text-secondary'>
		  	{themes.find(t => t.id == theme).icon}
		  </Dropdown.Toggle>

		  <Dropdown.Menu>
		  	{themes.map(t => (
				<Dropdown.Item active={t.id == theme} className='btn btn-sm' onClick={() => setTheme(t.id)}>
					{t.icon} {t.label}
				</Dropdown.Item>
			))}
		  </Dropdown.Menu>
		</Dropdown>
	)
}

export default ThemeSelector
