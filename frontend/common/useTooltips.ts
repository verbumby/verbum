import type * as TooltipModule from 'bootstrap/js/src/tooltip'
import { useEffect, useRef, useState } from 'react'

export function useTooltips<T extends HTMLElement>() {
    const el = useRef<T | null>(null)
    const [tooltipModule, setTooltipModule] = useState<
        typeof TooltipModule | null
    >(null)
    useEffect(() => {
        import('bootstrap/js/src/tooltip').then(setTooltipModule)
    }, [])

    useEffect(() => {
        if (!tooltipModule || !el.current) {
            return
        }

        const Tooltip = tooltipModule.default
        const ts: InstanceType<typeof Tooltip>[] = []
        for (const e of el.current.querySelectorAll(
            '[data-bs-toggle="tooltip"]',
        )) {
            ts.push(new Tooltip(e))
        }

        if (el.current.hasAttribute('data-bs-toggle')) {
            ts.push(new Tooltip(el.current))
        }

        return () => {
            for (const t of ts) {
                t.dispose()
            }
        }
    }, [tooltipModule, el])

    return el
}
