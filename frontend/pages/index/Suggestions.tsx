import * as React from 'react'
import { Suggestion } from '../../common'

type SuggestionsProps = {
    suggestions: Suggestion[],
    onClick: (s: string) => void,
}
export const Suggestions: React.VFC<SuggestionsProps> = ({ suggestions, onClick }) => (
    <ul className="suggestions">
        {suggestions.map(s => <li onClick={() => onClick(s)} key={s}>{s}</li>)}
    </ul>
)
