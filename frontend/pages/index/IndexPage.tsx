import * as React from 'react'
import { useEffect } from 'react'
import { useDispatch } from 'react-redux'
import { useHistory } from 'react-router-dom'
import { useDicts, useSearchState } from '../../store'
import { search, searchReset } from './search'
import { SearchControl } from './SearchControl'
import { DictsList } from './DictsList'
import { ArticleView, useURLSearch } from '../../common'

export const IndexPage: React.VFC = () => {
    const history = useHistory()
    const urlSearch = useURLSearch()
    const q: string = urlSearch.get('q') || ''
    const dicts = useDicts()
    const searchState = useSearchState()

    const dispatch = useDispatch()
    useEffect(() => {
        dispatch(search(urlSearch))
        return () => dispatch(searchReset())
    }, [q])

    const onSearch = (q: string) => {
        if (!q) {
            history.push('/')
        } else {
            history.push('/?q=' + encodeURIComponent(q))
        }
    }

    return (
        <div>
            <SearchControl onSearch={onSearch} urlQ={q} />
            <p />
            {q && searchState.hits.map(hit => <ArticleView key={`${hit.DictionaryID}-${hit.ID}`} a={hit} />)}
            {!q && <DictsList dictionaries={dicts} />}
        </div>
    )
}
