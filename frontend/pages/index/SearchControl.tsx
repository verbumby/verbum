import * as React from 'react'
import { useEffect, useState, useRef } from 'react'

import { Suggestions } from './Suggestions'
import { IconBackspace, IconSearch } from '../../icons'
import { Suggestion, useDelayed } from '../../common'

type SearchControlProps = {
    urlQ: string
    onSearch: (q: string) => void
}

export const SearchControl: React.VFC<SearchControlProps> = ({ urlQ, onSearch }) => {
    const [q, setQ] = useState<string>(urlQ)
    const qEl = useRef<HTMLInputElement>(null)

    let [
        suggestions,
        resetSuggestions,
        calculateQ,
        inputProps,
        suggestionViewProps,
    ] = useSuggestions()

    const urlQJustChanged = useRef<boolean>(false)
    useEffect(() => {
        setQ(urlQ)
        resetSuggestions()
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
        resetSuggestions()
        qEl.current.focus()
    }

    const { onChange } = inputProps
    inputProps = { ...inputProps, onChange: (e: React.ChangeEvent<HTMLInputElement>) => {
        setQ(e.target.value)
        onChange(e)
    }}

    return (
        <div id="search">
            <form action="/" method="get" onSubmit={e => { e.preventDefault(); onSearch(calculateQ(q)) }} >
                <div className="search-input">
                    <input
                        ref={qEl}
                        type="text"
                        name="q"
                        value={calculateQ(q)}
                        autoComplete="off"
                        {...inputProps}
                    />
                    {q && (<span className="btn button-clear" onClick={onClearClick}>
                        <IconBackspace />
                    </span>)}
                    <button type="submit" className="btn button-search button-search-wide">Шукаць</button>
                    <button type="submit" className="btn button-search button-search-small">
                        <IconSearch />
                    </button>
                </div>
                {suggestions.length > 0 && (
                    <Suggestions onClick={onSearch} {...suggestionViewProps} />
                )}
            </form>
        </div>
    )
}

type useSuggestionsSuggestionsViewProps = {
    suggestions: Suggestion[],
    active: number,
    setActive: (n: number) => void,
}

type useSuggestionsInputProps = {
    onChange: (e: React.ChangeEvent<HTMLInputElement>) => void,
    onKeyDown: (e: React.KeyboardEvent<HTMLInputElement>) => void,
    onBlur: (e: React.FocusEvent<HTMLInputElement>) => void
}

function useSuggestions(): [
    Suggestion[],
    () => void,
    (q: string) => string,
    useSuggestionsInputProps,
    useSuggestionsSuggestionsViewProps,
] {
    const [suggs, setSuggs] = useState<Suggestion[]>([])
    const [active, setActive] = useState<number>(-1)
    const [hard, setHard] = useState<boolean>(false)
    const promise = useRef<Promise<Suggestion[]>>(null)

    const resetSuggestions = () => {
        setSuggs([])
        setActive(-1)
        setHard(false)
    }

    const onChange = (e: React.ChangeEvent<HTMLInputElement>): void => {
        if (!e.target.value) {
            resetSuggestions()
            promise.current = null
        } else {
            setHard(false)
            // TODO: cancel prev request
            const p = verbumClient.suggest(e.target.value)
            promise.current = p
            p.then(suggs => {
                if (promise.current != p) {
                    return
                }
                if (active > suggs.length - 1) {
                    setActive(-1)
                    setHard(false)
                }
                setSuggs(suggs)
            })
        }
    }

    const onKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
        switch (e.key) {
            case "Escape":
                resetSuggestions()
                break
            case "ArrowDown":
                if (active + 1 < suggs.length) {
                    setActive(active + 1)
                    setHard(true)
                } else {
                    setActive(-1)
                    setHard(true)
                }
                break
            case "ArrowUp":
                if (active - 1 >= -1) {
                    setActive(active - 1)
                    setHard(true)
                } else {
                    setActive(suggs.length - 1)
                    setHard(true)
                }
                break
        }
    }

    const onBlur = (e: React.FocusEvent<HTMLInputElement>) => {
        if (suggs.length > 0) {
            setTimeout(() => resetSuggestions(), 150)
        }
    }

    const setActiveSuggestionDelayed = useDelayed((n: number) => {
        setActive(n)
        setHard(false)
    }, 15)

    const calculateQ = (q: string): string => {
        return hard && active > -1 ? suggs[active] : q
    }

    return [
        suggs,
        resetSuggestions,
        calculateQ,
        {onChange, onKeyDown, onBlur},
        {suggestions: suggs, active, setActive: setActiveSuggestionDelayed},
    ]
}