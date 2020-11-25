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

    const triggerSearch = () => {
        
    }

    return (
        <div id="search">
            <form action="/" method="get" onSubmit={e => { e.preventDefault(); triggerSearch() } } >
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
                    <button type="submit" className="btn button-search button-search-wide">Шукаць</button>
                    <button type="submit" className="btn button-search button-search-small">
                        <IconSearch />
                    </button>
                </div>
                <ul className="suggestions" style={{ display: 'none' }}>
                </ul>
            </form>
        </div>
    )
}
