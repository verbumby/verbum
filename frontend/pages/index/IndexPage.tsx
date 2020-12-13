import * as React from 'react'
import { useEffect } from 'react'
import { useDispatch } from 'react-redux'
import { useRouteMatch } from 'react-router-dom'
import { Helmet } from "react-helmet"

import { useDicts, useSearchState } from '../../store'
import { search, searchReset } from './search'
import { DictsList } from './DictsList'
import { ArticleView, useURLSearch, SearchControl } from '../../common'

export const IndexPage: React.VFC = () => {
    const match = useRouteMatch()
    const urlSearch = useURLSearch()
    const q: string = urlSearch.get('q') || ''
    const dicts = useDicts()
    const searchState = useSearchState()

    const dispatch = useDispatch()
    useEffect(() => {
        dispatch(search(match, urlSearch))
        return () => dispatch(searchReset())
    }, [q])

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
            <SearchControl urlQ={q} />
            {!q ? renderDictList() : renderSearchResults() }
        </div>
    )
}
