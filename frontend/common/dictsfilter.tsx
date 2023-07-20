import * as React from 'react'
import { useEffect, useState, useCallback } from 'react'

import { IconFunnel } from '../icons'
import { useListedDicts } from '../store'

type DictsFilterProps = {
    state: string
    onChange?: (state: string) => void
}

const DictsFilter: React.FC<DictsFilterProps> = ({ state, onChange }) => {
    const dicts = useListedDicts()
    let checkedDicts: string[]
    if (state == '') {
        checkedDicts = dicts.map(d => d.ID)
    } else if (state == '-') {
        checkedDicts = []
    } else {
        checkedDicts = state
            .split(',')
            .filter(id => dicts.findIndex(d => d.ID == id) >= 0)
    }

    const isChecked = (id: string): boolean => {
        return checkedDicts.includes(id)
    }

    const onInputChange = (id: string) => (e: React.ChangeEvent<HTMLInputElement>) => {
        const toEncode = [...checkedDicts]

        if (e.target.checked) {
            toEncode.push(id)
        } else {
            toEncode.splice(toEncode.indexOf(id), 1)
        }

        let newState: string
        if (toEncode.length == 0) {
            newState = '-'
        } else if (dicts.every(d => toEncode.includes(d.ID))) {
            newState = ''
        } else {
            newState = dicts.map(d => d.ID).filter(id => toEncode.includes(id)).join(',')
        }

        onChange(newState)
    }

    return (
        <div className="dicts-filter">
            <form>
                <div className="mb-2">
                    Шукаць <span className="btn btn-link" onClick={() => onChange('')}>усюды</span>
                    , <span className="btn btn-link" onClick={() => onChange('-')}>нідзе</span>
                    , у:
                </div>
                <ul>
                    {dicts.map(({ ID, Title }) =>
                        <li>
                            <div className="form-check">
                                <input
                                    className="form-check-input"
                                    type="checkbox"
                                    checked={isChecked(ID)}
                                    onChange={onInputChange(ID)}
                                    id={`dicts-filter-item-${ID}`}
                                />
                                <label className="form-check-label" htmlFor={`dicts-filter-item-${ID}`}>{Title}</label>
                            </div>
                        </li>
                    )}
                </ul>
            </form>
        </div>
    )
}

export function useDictsFilter(urlIn: string): {
    inDicts: string,
    icon: JSX.Element,
    filter: JSX.Element,
} {
    const [inDicts, setIn] = useState<string>(urlIn)
    const [shown, setShown] = useState<boolean>(false)

    useEffect(() => setIn(urlIn), [urlIn])

    let styles: React.CSSProperties = {}
    if (inDicts != '') {
        styles = { color: 'red' }
    }

    const windowClickCallback = useCallback((e: MouseEvent) => {
        let el = e.target as HTMLElement
        while (el.parentElement) {
            if (el.classList.contains('dicts-filter')) {
                return
            }
            el = el.parentElement
        }

        setShown(false)
        window.removeEventListener('click', windowClickCallback)
        window.removeEventListener('keydown', windowKeyDownCallback)
    }, [])

    const windowKeyDownCallback = useCallback((e: KeyboardEvent) => {
        if (e.key !== "Escape") {
            return
        }

        setShown(false)
        window.removeEventListener('click', windowClickCallback)
        window.removeEventListener('keydown', windowKeyDownCallback)
    }, [])

    const icon = (
        <span
            className="btn button-control button-funnel"
            style={styles}
            onClick={(e) => {
                e.stopPropagation()
                setShown(!shown)
                if (!shown) {
                    window.addEventListener('click', windowClickCallback)
                    window.addEventListener('keydown', windowKeyDownCallback)
                } else {
                    window.removeEventListener('click', windowClickCallback)
                    window.removeEventListener('keydown', windowKeyDownCallback)
                }
            }}
        >
            <IconFunnel fill={inDicts != ''} />
        </span>
    )

    const filter = shown ? <DictsFilter state={inDicts} onChange={setIn} /> : <></>

    return { inDicts, icon, filter }
}
