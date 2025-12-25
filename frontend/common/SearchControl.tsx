import * as React from 'react'
import { useEffect, useState, useRef, useCallback } from 'react'
import { Suggestions } from './Suggestions'
import { IconBackspace, IconSearch } from '../icons'
import { Dict, Suggestion, useDelayed, useDispatch } from '.'
import { useNavigate, To } from 'react-router'
import { hideLoading, showLoading } from 'react-redux-loading-bar'
import { useDictsFilter } from './dictsfilter'

type SearchControlProps = {
    inBound: Dict[]
    urlQ: string
    urlIn: string
    calculateSearchURL: (q: string, inDicts: string) => To
    filterEnabled: boolean
}

export const SearchControl: React.FC<SearchControlProps> = ({ inBound, urlQ, urlIn, calculateSearchURL, filterEnabled }) => {
    const [q, setQ] = useState<string>(urlQ)
    const qEl = useRef<HTMLInputElement>(null)
    const navigate = useNavigate()

    const {
        inDicts,
        icon: dictsFilterIcon,
        filter: dictsFilter
    } = useDictsFilter(inBound, urlIn)

    const onSearch = (q: string) => {
        navigate(calculateSearchURL(q, inDicts))
    }

    let [
        suggestions,
        resetSuggestions,
        calculateQ,
        inputProps,
        suggestionViewProps,
    ] = useSuggestions(inDicts || inBound.map(d => d.ID).join(','))

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
    inputProps = {
        ...inputProps, onChange: (e: React.ChangeEvent<HTMLInputElement>) => {
            setQ(e.target.value)
            onChange(e)
        }
    }


    useEffect(() => {
        const globalKeyPressHandler = (ev: KeyboardEvent) => {
            if (ev.key === '/' && ev.target !== qEl.current) {
                ev.preventDefault()
                ev.stopPropagation()
                window.setTimeout(() => {
                    qEl.current.scrollIntoView({ behavior: 'smooth', block: 'center', inline: "center" })
                    qEl.current.focus({ preventScroll: true })
                    qEl.current.setSelectionRange(0, qEl.current.value.length)
                }, 10)
            }
        }

        window.addEventListener('keypress', globalKeyPressHandler)
        return () => window.removeEventListener('keypress', globalKeyPressHandler)
    }, [])

    return (
        <div id="search">
            <form
                action="/"
                method="get"
                onSubmit={
                    e => {
                        e.preventDefault()
                        if (inDicts === '-') {
                            return
                        }
                        onSearch(calculateQ(q))
                    }
                }
            >
                <div className="search-input">
                    <input
                        ref={qEl}
                        type="text"
                        name="q"
                        value={calculateQ(q)}
                        autoComplete="off"
                        {...inputProps}
                    />
                    {q && (<span className="btn button-control button-clear" onClick={onClearClick}>
                        <IconBackspace />
                    </span>)}
                    {filterEnabled ? dictsFilterIcon : null}
                    <button type="submit" disabled={inDicts === '-'} className="btn button-search button-search-wide">Шукаць</button>
                    <button type="submit" disabled={inDicts === '-'} className="btn button-search button-search-small">
                        <IconSearch />
                    </button>
                </div>
                {suggestions.length > 0 && (
                    <Suggestions onClick={onSearch} {...suggestionViewProps} />
                )}
                {filterEnabled ? dictsFilter : null}
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
}

function useSuggestions(inDicts: string): [
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
    const abort = useRef<AbortController>(null)
    const dispatch = useDispatch()

    const resetSuggestions = useCallback(() => {
        setSuggs([])
        setActive(-1)
        setHard(false)
        promise.current = null
        if (abort.current) {
            abort.current.abort()
        }
        abort.current = null
        // onChangeHandlerCancel()
        window.removeEventListener('click', onWindowClick)
    }, [])

    const [onChangeHandler, onChangeHandlerCancel] = useDelayed((q: string): void => {
        if (!q || inDicts == '-') {
            resetSuggestions()
        } else {
            setHard(false)

            dispatch(showLoading())

            if (abort.current) {
                abort.current.abort()
            }
            abort.current = new AbortController()

            const p = verbumClient
                .withSignal(abort.current.signal)
                .suggest(q, inDicts)
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
                if (suggs.length > 0) {
                    window.addEventListener('click', onWindowClick)
                } else {
                    window.removeEventListener('click', onWindowClick)
                }
            }).catch(() => {
                // ignore abort exception
            }).finally(() => {
                dispatch(hideLoading())
            })
        }
    }, 150)

    const onChange = (e: React.ChangeEvent<HTMLInputElement>): void => {
        onChangeHandler(e.target.value)
    }

    const onKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
        if (e.key == "Escape") {
            resetSuggestions()
        } else if (e.key == "ArrowDown" || e.key == "j" && e.metaKey) {
            e.stopPropagation()
            e.preventDefault()
            if (active + 1 < suggs.length) {
                setActive(active + 1)
                setHard(true)
            } else {
                setActive(-1)
                setHard(true)
            }
        } else if (e.key == "ArrowUp" || e.key == "k" && e.metaKey) {
            e.stopPropagation()
            e.preventDefault()
            if (active - 1 >= -1) {
                setActive(active - 1)
                setHard(true)
            } else {
                setActive(suggs.length - 1)
                setHard(true)
            }
        }
    }

    const onWindowClick = useCallback(() => {
        resetSuggestions()
    }, [])

    const [setActiveSuggestionDelayed] = useDelayed((n: number) => {
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
        { onChange, onKeyDown },
        { suggestions: suggs, active, setActive: setActiveSuggestionDelayed },
    ]
}
