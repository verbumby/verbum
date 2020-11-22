import * as React from 'react'

import { IconBackspace } from './IconBackspace'
import { IconSearch } from './IconSearch'

export const SearchControl: React.VFC = () => {
    const [q, setQ] = React.useState<string>('')
    const qEl = React.useRef<HTMLInputElement>(null)

    const onClearClick = () => {
        setQ('')
        qEl.current.focus()
    }

    return (
        <div id="search">
            <form action="/" method="get">
                <div className="search-input">
                    <input
                        ref={qEl}
                        type="text"
                        name="q"
                        value={q}
                        onChange={(e) => setQ(e.target.value)}
                        autoComplete="off"
                        autoFocus
                    />
                    {q && (<span className="btn button-clear" onClick={onClearClick}>
                        <IconBackspace />
                    </span>)}
                    <span className="btn button-search button-search-wide">Шукаць</span>
                    <span className="btn button-search button-search-small">
                        <IconSearch />
                    </span>
                </div>
                <ul className="suggestions" style={{ display: 'none' }}>
                </ul>
            </form>
        </div>
    )
}
