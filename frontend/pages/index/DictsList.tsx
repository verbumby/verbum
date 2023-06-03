import * as React from 'react'
import { Link } from 'react-router-dom'
import { Dict } from '../../common'

type DictsListProps = {
    dictionaries: Dict[]
}

export const DictsList: React.FC<DictsListProps> = ({ dictionaries }) => {
    return (
        <ul>
            {dictionaries.map(d => <li key={d.ID}><Link to={`/${d.ID}`}>{d.Title}</Link></li>)}
        </ul>
    )
}
