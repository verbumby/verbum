import * as React from 'react'
import { Helmet } from 'react-helmet'
import { useDispatch } from 'react-redux'
import { useRouteMatch } from 'react-router-dom'
import { ArticleView, PaginationView, SearchControl } from '../../common'
import { useDict, useDictArticles, useLetterFilter } from '../../store'
import { letterFilterFetch, letterFilterReset } from './letterfilter'
import { dictArticlesFetch, MatchParams, dictArticlesReset, useURLSearch } from './dict'
import { LetterFilterView } from './LetterFilterView'
import { useURLSearch as useIndexURLSearch } from '../index/search'

export const DictPage: React.VFC = ({ }) => {
    const match = useRouteMatch<MatchParams>()
    const urlSearch = useURLSearch()

    const prefix = urlSearch.get('prefix')
    const page = urlSearch.get('page')

    const dict = useDict(match.params.dictID)
    const letterFilter = useLetterFilter()
    const dictArticles = useDictArticles()
    const dispatch = useDispatch()

    const indexURLSearch = useIndexURLSearch()

    React.useEffect(() => {
        dispatch(letterFilterFetch(match, urlSearch))
    }, [match.params.dictID, prefix])
    React.useEffect(() => () => dispatch(letterFilterReset()), [])

    React.useEffect(() => {
        dispatch(dictArticlesFetch(match, urlSearch))
        return () => dispatch(dictArticlesReset())
    }, [match.params.dictID, prefix, page])

    if (!letterFilter) {
        return <></>
    }

    return (
        <>
            <Helmet>
                <title>{dict.Title}</title>
                <meta name="description" content={`${dict.Title}`} />
                <meta property="og:title" content={`${dict.Title}`} />
                <meta property="og:description" content={`${dict.Title}`} />
                <meta name="robots" content="noindex, follow" />
            </Helmet>
            <h4 className="ml-1 mr-1 mb-3">{dict.Title}</h4>
            <div className="mb-3">
                <SearchControl
                    urlQ=""
                    urlIn={dict.ID}
                    calculateSearchURL={
                        (q, inDicts) => indexURLSearch
                            .clone()
                            .set('q', q)
                            .set('in', inDicts)
                            .set('page', 1)
                            .encode()
                    }
                />
            </div>
            <LetterFilterView
                letterFilter={letterFilter}
                prefixToURL={prefix => ({
                    search: urlSearch
                        .clone()
                        .reset('page')
                        .set('prefix', prefix)
                        .encode()
                })}
            />
            {dictArticles && <>
                <div>
                    {dictArticles.Articles.map(a => <ArticleView key={a.DictionaryID + a.ID} a={a} />)}
                </div>
                <PaginationView
                    current={dictArticles.Pagination.Current}
                    total={dictArticles.Pagination.Total}
                    pageToURL={p => ({
                        search: urlSearch
                            .clone()
                            .set('page', p)
                            .encode()
                    })}
                />
            </>}
        </>
    )
}
