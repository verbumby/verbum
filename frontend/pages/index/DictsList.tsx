import * as React from 'react'
import { Link } from 'react-router'
import { Dict, DictTitle } from '../../common'

type DictsListProps = {
    dictionaries: Dict[]
}

export const DictsList: React.FC<DictsListProps> = ({ dictionaries }) => {
    return (
        <ul className='mt-2'>
            {dictionaries.map(d => <li key={d.ID}><Link to={`/${d.ID}`}><DictTitle d={d} /></Link></li>)}
        </ul>
    )
}
