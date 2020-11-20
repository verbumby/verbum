import * as React from 'react'
import { useEffect } from 'react'
import { useDispatch } from 'react-redux'
import { DictionariesListState, useSelector, dictionariesListFetch } from '../reducers'

export const IndexPage: React.VFC = () => {
    const dicts = useSelector<DictionariesListState>(state => state.dictionaries)
    const dispatch = useDispatch()
    useEffect(() => {
        if (dicts.length == 0) {
            dispatch(dictionariesListFetch())
        }
    })

    return (<p>
        <ul>
            {dicts.map(d => <li key={d.ID}><a href={`/${d.ID}`}>{d.Title}</a></li>)}
        </ul>
    </p>)
}
