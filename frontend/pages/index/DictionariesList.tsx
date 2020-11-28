import * as React from 'react'
import { Link } from 'react-router-dom'
import { Dictionary } from '../../reducers'

type DictionariesListProps = {
    dictionaries: Dictionary[]
}

export const DictionariesList: React.VFC<DictionariesListProps> = ({ dictionaries }) => {
    return (
        <ul>
            {dictionaries.map(d => <li key={d.ID}><Link to={`/${d.ID}`}>{d.Title}</Link></li>)}
        </ul>
    )
}
