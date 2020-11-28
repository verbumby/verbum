import * as React from 'react'
import { useEffect } from 'react'
import { useDispatch } from 'react-redux'
import { useLocation } from 'react-router-dom'
import {
    useDicts,
    search,
    searchReset,
    useSearchState,
} from '../../reducers'
import { SearchControl } from './SearchControl'
import { DictionariesList } from './DictionariesList'
import { Article } from '../../common/Article'

export const IndexPage: React.VFC = () => {
    const urlSearch = useURLSearchQuery()
    const q: string = urlSearch.get('q') || ''
    const dicts = useDicts()
    const searchState = useSearchState()

    const dispatch = useDispatch()
    useEffect(() => {
        if (!q) {
            return
        }
        dispatch(search(q))
        return () => dispatch(searchReset())
    }, [q])

    return (
        <div>
            <SearchControl />
            <p />
            {q && searchState.hits.map(hit => <Article key={`${hit.DictionaryID}-${hit.ID}`} a={hit} />)}
            {!q && <DictionariesList dictionaries={dicts} />}
        </div>
    )
}

function useURLSearchQuery(): URLSearchParams {
    return new URLSearchParams(useLocation().search)
}
