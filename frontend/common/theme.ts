export type ThemeID = 'light'|'dark'|'auto'

export type Theme = {
	id: ThemeID
	icon: React.ReactElement,
	label: string,
}
