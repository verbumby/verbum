import type * as Bootstrap from 'bootstrap'
import { useEffect, useRef, useState } from 'react'

export function useTooltips<T extends HTMLElement>() {
    const el = useRef<T | null>(null)
    const [bootstrapAPI, setBootstrapAPI] = useState<typeof Bootstrap | null>(null)
    useEffect(() => {
        import('bootstrap').then(setBootstrapAPI)
    }, [])

    useEffect(() => {
        if (!bootstrapAPI || !el.current) {
            return
        }

        const ts: InstanceType<typeof Bootstrap.Tooltip>[] = []
        for (const e of el.current.querySelectorAll(
            '[data-bs-toggle="tooltip"]',
        )) {
            ts.push(new bootstrapAPI.Tooltip(e))
        }

        if (el.current.hasAttribute('data-bs-toggle')) {
            ts.push(new bootstrapAPI.Tooltip(el.current))
        }

        return () => {
            for (const t of ts) {
                t.dispose()
            }
        }
    }, [bootstrapAPI, el])

    return el
}
