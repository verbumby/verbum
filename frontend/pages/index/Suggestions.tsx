import * as React from 'react'
import { Suggestion } from '../../common'

type SuggestionsProps = {
    suggestions: Suggestion[],
}
export const Suggestions: React.VFC<SuggestionsProps> = ({ suggestions }) => (
    <ul className="suggestions">
        {suggestions.map(s => <li key={s}>{s}</li>)}
    </ul>
)
