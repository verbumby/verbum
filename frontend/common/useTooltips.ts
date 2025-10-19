import { useEffect, useRef, useState } from "react";

export function useTooltips<T extends HTMLElement>() {
	const el = useRef<T>()
    const [bootstrapAPI, setBootstrapAPI] = useState(null)
    useEffect(() => {
        import('bootstrap').then(setBootstrapAPI)
    }, [])

    useEffect(() => {
        if (!bootstrapAPI || !el.current) {
            return
        }

        let ts = new Array()
        for (let e of el.current.querySelectorAll('[data-bs-toggle="tooltip"]')) {
            ts.push(new bootstrapAPI.Tooltip(e))
        }

        return () => {
            for (let t of ts) {
                t.dispose()
            }
        }
    }, [bootstrapAPI, el])

	return el
}
