import * as React from 'react'
import { Suggestion } from '../../common'

type SuggestionsProps = {
    suggestions: Suggestion[],
    activeOne: string,
    onClick: (s: string) => void,
    setActiveOne: (s: string) => void,
}
export const Suggestions: React.VFC<SuggestionsProps> = ({ suggestions, activeOne, onClick, setActiveOne }) => (
    <ul className="suggestions">
        {suggestions.map(s => (
            <li
                key={s}
                className={s === activeOne ? 'active': ''}
                onClick={() => onClick(s)}
                onMouseEnter={() => setActiveOne(s)}
                onMouseLeave={() => setActiveOne('')}
            >{s}</li>
        ))}
    </ul>
)
