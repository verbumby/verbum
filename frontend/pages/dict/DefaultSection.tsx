import * as React from 'react'
import { FC, useEffect } from 'react'

import { Helmet } from 'react-helmet'
import { Link, useRouteMatch } from 'react-router-dom'
import { ArticleView, NoSearchResults, PaginationView, SearchControl, useDispatch } from '../../common'
import { useDict, useDictArticles, useLetterFilter } from '../../store'
import { letterFilterFetch, letterFilterReset } from './letterfilter'
import { dictArticlesFetch, MatchParams, dictArticlesReset, useURLSearch } from './dict'
import { LetterFilterView } from './LetterFilterView'

export const DefaultSection: FC = ({ }) => {
    const match = useRouteMatch<MatchParams>()
    const urlSearch = useURLSearch()

    const [dict] = useDict(match.params.dictID)

    const letterFilter = useLetterFilter()
    const dictArticles = useDictArticles()
    const dispatch = useDispatch()

    const q = urlSearch.get('q')
    const prefix = urlSearch.get('prefix')
    const page = urlSearch.get('page')

    useEffect(() => {
        dispatch(letterFilterFetch(match, urlSearch))
    }, [match.params.dictID, prefix])
    useEffect(() => () => { dispatch(letterFilterReset()) }, [])

    useEffect(() => {
        dispatch(dictArticlesFetch(match, urlSearch))
        return () => { dispatch(dictArticlesReset()) }
    }, [match.params.dictID, prefix, page, q])

    if (!letterFilter) {
        return <></>
    }

    const topLinks: React.JSX.Element[] = []
    const pushToTopLinks = (node: React.JSX.Element, key: string) => {
        if (topLinks.length > 0) {
            topLinks.push(<span
                key={`separator-${key}`}
                style={{ color: 'darkgray' }}> ∙ </span>)
        }
        topLinks.push(React.cloneElement(node, { key }))
    }

    if (dict.HasPreface) {
        pushToTopLinks(<Link to={{
            pathname: match.url,
            search: urlSearch.clone()
                .resetAll()
                .set('section', 'preface')
                .encode()
        }}>Прадмова</Link>, 'preface')
    }

    if (dict.HasAbbrevs) {
        pushToTopLinks(<Link to={{
            pathname: match.url,
            search: urlSearch.clone()
                .resetAll()
                .set('section', 'abbr')
                .encode()
        }}>Скарачэнні</Link>, 'abbr')
    }

    if (dict.ScanURL) {
        pushToTopLinks(<a href={dict.ScanURL} target='_blank' rel='noopener noreferrer'>Кніга ў PDF/DjVu</a>, 'scan')
    }

    return (
        <>
            <Helmet>
                <title>{dict.Title}</title>
                <meta name="description" content={dict.Title} />
                <meta property="og:title" content={dict.Title} />
                <meta property="og:description" content={dict.Title} />
                <meta name="robots" content={`${urlSearch.allDefault() ? '' : 'no'}index, follow`} />
            </Helmet>
            <div className='mx-1 mb-3'>
                <h4>{dict.Title}</h4>
                {topLinks}
            </div>
            <div className="mb-3">
                <SearchControl
                    inBound={[dict]}
                    urlQ={q}
                    urlIn={dict.ID}
                    filterEnabled={false}
                    calculateSearchURL={
                        (q, inDicts) => ({
                            search: urlSearch
                                .clone()
                                .set('q', q)
                                .set('page', 1)
                                .encode()
                        })
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
            {dictArticles && dictArticles.Articles.length > 0 && <>
                <div>
                    {dictArticles.Articles.map(a => (
                        <ArticleView
                            key={a.DictionaryID + a.ID}
                            a={a}
                            showExternalButton={true}
                            showSource={false}
                        />
                    ))}
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
            {dictArticles && dictArticles.Articles.length == 0
                && <NoSearchResults q={q} suggestions={dictArticles.TermSuggestions}
                    calculateSuggestionURL={s => ({ search: urlSearch.clone().set('q', s).encode() })} />}
        </>
    )
}
