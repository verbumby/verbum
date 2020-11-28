import * as React from 'react'
import { useEffect } from 'react'
import { useHistory } from 'react-router-dom'
import { useLocation } from 'react-router-dom'

import { IconBackspace, IconSearch } from '../../icons'

export const SearchControl: React.VFC = () => {
    const history = useHistory()
    const urlQ = useURLSearchQuery().get('q') || ''
    const [q, setQ] = React.useState<string>(urlQ)
    const qEl = React.useRef<HTMLInputElement>(null)

    useEffect(() => {
        qEl.current.setSelectionRange(0, qEl.current.value.length)
    }, [urlQ])

    const onClearClick = () => {
        setQ('')
        qEl.current.focus()
    }

    const triggerSearch = () => {
        history.push('/?q=' + encodeURIComponent(q))
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

// TODO: deduplicate this and move to some utils
function useURLSearchQuery(): URLSearchParams {
    return new URLSearchParams(useLocation().search)
}
