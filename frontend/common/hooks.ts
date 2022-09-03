import { useLocation } from 'react-router-dom'
import { useRef } from 'react'
import { URLSearch, URLSearchEntries } from './urlsearch'

export function useURLSearch<Entries extends URLSearchEntries>(defaults: Entries): URLSearch<Entries> {
    return new URLSearch<Entries>(defaults, new URLSearchParams(useLocation().search))
}

export function useDelayed<T extends (...args: any[]) => any>(f: T, delayMs: number): [(...funcArgs: Parameters<T>) => void, () => void] {
    const timeoutID = useRef<number>(null)
    return [(...args: Parameters<T>): void => {
        if (timeoutID.current) {
            window.clearTimeout(timeoutID.current)
        }
        timeoutID.current = window.setTimeout(() => {
            f(...args)
            timeoutID.current = null
        }, delayMs)
    }, (): void => {
        if (timeoutID.current) {
            window.clearTimeout(timeoutID.current)
        }
        timeoutID.current = null
    }]
}
