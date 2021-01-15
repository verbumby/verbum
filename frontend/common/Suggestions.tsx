import * as React from 'react'
import { Suggestion } from '.'

type SuggestionsProps = {
    suggestions: Suggestion[],
    active: number,
    onClick: (s: string) => void,
    setActive: (n: number) => void,
}
export const Suggestions: React.VFC<SuggestionsProps> = ({ suggestions, active, onClick, setActive }) => (
    <ul className="suggestions">
        {suggestions.map((s, i) => (
            <li
                key={s}
                className={i === active ? 'active': ''}
                onClick={() => onClick(s)}
                onMouseEnter={() => setActive(i)}
                onMouseLeave={() => setActive(-1)}
            >{s}</li>
        ))}
    </ul>
)
