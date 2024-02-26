import { ThemeID } from "./common"

export {}

declare global {
    interface Window {
        updateStoredTheme: (theme: ThemeID) => void
        getStoredTheme: () => ThemeID
    }
}

window.getStoredTheme = () => {
    const stored = localStorage.getItem('theme')
    return stored ? stored as ThemeID : 'auto'
}

window.updateStoredTheme = (theme) => {
    if (theme == 'auto') {
        localStorage.removeItem('theme')
    } else {
        localStorage.setItem('theme', theme)
    }
    setTheme(theme)
}

window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
    const storedTheme = window.getStoredTheme()
    if (storedTheme !== 'light' && storedTheme !== 'dark') {
        setTheme(storedTheme)
    }
})

const setTheme = (theme: ThemeID) => {
    if (theme === 'auto') {
        document.documentElement.setAttribute('data-bs-theme', (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'))
    } else {
        document.documentElement.setAttribute('data-bs-theme', theme)
    }
}

setTheme(window.getStoredTheme())
