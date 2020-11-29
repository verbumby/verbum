import * as React from 'react'
import { useEffect } from 'react'

import { IconBackspace, IconSearch } from '../../icons'

type SearchControlProps = {
    urlQ: string
    onSearch: (q: string) => void
}

export const SearchControl: React.VFC<SearchControlProps> = ({ urlQ, onSearch }) => {
    const [q, setQ] = React.useState<string>(urlQ)
    const qEl = React.useRef<HTMLInputElement>(null)

    useEffect(() => {
        setQ(urlQ)
        qEl.current.setSelectionRange(0, qEl.current.value.length)
    }, [urlQ])

    const onClearClick = () => {
        setQ('')
        qEl.current.focus()
    }

    return (
        <div id="search">
            <form action="/" method="get" onSubmit={e => { e.preventDefault(); onSearch(q) } } >
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
