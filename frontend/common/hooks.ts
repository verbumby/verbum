import { useLocation } from 'react-router-dom'
import { useRef } from 'react'

export function useURLSearch(): URLSearchParams {
    return new URLSearchParams(useLocation().search)
}

export function useDelayed<T extends (...args: any[]) => any>(f: T, delayMs: number): (...funcArgs: Parameters<T>) => void {
    const timeoutID = useRef<number>(null)
    return (...args: Parameters<T>): void => {
        if (timeoutID.current) {
            window.clearTimeout(timeoutID.current)
        }
        timeoutID.current = window.setTimeout(() => {
            f(...args)
            timeoutID.current = null
        }, delayMs)
    }
}
