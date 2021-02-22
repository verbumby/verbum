import * as React from 'react'
import { useEffect } from 'react'
import { useDispatch } from 'react-redux'
import { Link, useRouteMatch } from 'react-router-dom'
import { Helmet } from "react-helmet"

import { useDicts, useSearchState } from '../../store'
import { search, searchReset, useURLSearch } from './search'
import { DictsList } from './DictsList'
import { ArticleView, SearchControl } from '../../common'

export const IndexPage: React.VFC = () => {
    const match = useRouteMatch()
    const urlSearch = useURLSearch()
    const q = urlSearch.get('q')
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
                <meta name="description" content="Verbum - Анлайн Слоўнік Беларускай Мовы" />
                <meta property="og:title" content="Verbum - Анлайн Слоўнік Беларускай Мовы" />
                <meta property="og:description" content="Verbum - Анлайн Слоўнік Беларускай Мовы" />
                <meta name="robots" content="index, follow" />
            </Helmet>
            <p />
            <DictsList dictionaries={dicts} />
        </>
    )

    const renderSearchResults = (): React.ReactNode => {


        return (
            <>
                <Helmet>
                    <title>{q} - Пошук</title>
                    <meta name="description" content={`${q} - Пошук`} />
                    <meta property="og:title" content={`${q} - Пошук`} />
                    <meta property="og:description" content={`${q} - Пошук`} />
                    <meta name="robots" content="noindex, nofollow" />
                </Helmet>
                {searchState.searchResult && searchState.searchResult.Articles.map(
                    hit => <ArticleView key={`${hit.DictionaryID}-${hit.ID}`} a={hit} />
                )}
                {searchState.searchResult
                    && searchState.searchResult.Articles.length == 0
                    && renderNoSearchResults()
                }
            </>
        )
    }

    const renderNoSearchResults = (): React.ReactNode => {
        const suggs = searchState.searchResult.TermSuggestions
        return (
            <div className="no-results">
            <p>Па запыце <strong>{q}</strong> нічога не знойдзена.</p>
            {suggs.length == 1 && (
                <p>Магчыма вы шукалі&nbsp;
                    <Link to={{search: urlSearch.clone().set('q', suggs[0]).encode()}}>
                        {suggs[0]}
                    </Link>.
                </p>
            )}
            {suggs.length > 1 && (
                <p>Магчыма вы шукалі:
                    <ul>
                        {suggs.map(s => (
                            <li key={s}>
                                <Link to={{search: urlSearch.clone().set('q', s).encode()}}>
                                    {s}
                                </Link>
                            </li>
                        ))}
                    </ul>
                </p>
            )}
        </div>
        )
    }

    return (
        <div>
            <SearchControl urlQ={q} />
            {!q ? renderDictList() : renderSearchResults()}
        </div>
    )
}
