import * as React from 'react'
import { useEffect } from 'react'
import { useDispatch } from 'react-redux'
import { useHistory } from 'react-router-dom'
import { Helmet } from "react-helmet"

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

    const renderDictList = (): React.ReactNode => (
        <>
            <Helmet>
                <title>Verbum - Анлайн Слоўнік Беларускай Мовы</title>
                <meta name="robots" content="index, follow" />
            </Helmet>
            <p />
            <DictsList dictionaries={dicts} />
        </>
    )

    const renderSearchResults = (): React.ReactNode => (
        <>
            <Helmet>
                <title>face - Пошук</title>
                <meta name="robots" content="noindex, nofollow" />
            </Helmet>
            {searchState.hits.map(hit => <ArticleView key={`${hit.DictionaryID}-${hit.ID}`} a={hit} />)}
        </>
    )

    return (
        <div>
            <SearchControl onSearch={onSearch} urlQ={q} />
            {!q ? renderDictList() : renderSearchResults() }
        </div>
    )
}
