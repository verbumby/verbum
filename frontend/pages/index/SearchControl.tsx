import * as React from 'react'
import { useEffect, useState, useRef } from 'react'

import { Suggestions } from './Suggestions'
import { IconBackspace, IconSearch } from '../../icons'
import { Suggestion } from '../../common'

type SearchControlProps = {
    urlQ: string
    onSearch: (q: string) => void
}

export const SearchControl: React.VFC<SearchControlProps> = ({ urlQ, onSearch }) => {
    const [q, setQ] = useState<string>(urlQ)
    const qEl = useRef<HTMLInputElement>(null)

    const urlQJustChanged = useRef<boolean>(false)
    useEffect(() => {
        setQ(urlQ)
        setSuggestions([])
        urlQJustChanged.current = true
    }, [urlQ])
    useEffect(() => {
        if (urlQJustChanged.current) {
            qEl.current.focus()
            qEl.current.setSelectionRange(0, qEl.current.value.length)
            urlQJustChanged.current = false
        }
    }, [q, urlQJustChanged.current])

    const onClearClick = () => {
        setQ('')
        qEl.current.focus()
    }

    const [suggestions, setSuggestions] = useState<Suggestion[]>([])

    const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const v = e.target.value
        setQ(v)
        if (!v) {
            setSuggestions([])
        } else {
            // TODO: cancel prev request and check if it's the same promise
            verbumClient.suggest(v).then(suggs => setSuggestions(suggs))
        }
    }

    const onBlur = (e: React.FocusEvent<HTMLInputElement>) => {
        if (suggestions.length > 0) {
            setTimeout(() => setSuggestions([]), 150)
        }
    }

    return (
        <div id="search">
            <form action="/" method="get" onSubmit={e => { e.preventDefault(); onSearch(q) }} >
                <div className="search-input">
                    <input
                        ref={qEl}
                        type="text"
                        name="q"
                        value={q}
                        onChange={onChange}
                        onBlur={onBlur}
                        autoComplete="off"
                    />
                    {q && (<span className="btn button-clear" onClick={onClearClick}>
                        <IconBackspace />
                    </span>)}
                    <button type="submit" className="btn button-search button-search-wide">Шукаць</button>
                    <button type="submit" className="btn button-search button-search-small">
                        <IconSearch />
                    </button>
                </div>
                {suggestions.length > 0 && <Suggestions onClick={onSearch} suggestions={suggestions} />}
            </form>
        </div>
    )
}
