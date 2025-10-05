import * as React from 'react'
import { useEffect } from 'react'
import { NoSearchResults, NotFound, useDispatch } from '../../common'
import { Link, useParams } from 'react-router'
import { Helmet } from "react-helmet"

import { useDictsInSection, useSearchState, useSection, useSections } from '../../store'
import { search, searchReset, useURLSearch } from './search'
import { DictsList } from './DictsList'
import { ArticleView, PaginationView, SearchControl } from '../../common'

export const IndexPage: React.FC = () => {
    const params = useParams<{ sectionID?: string }>()
    const sectionID = params.sectionID || 'default'
    const section = useSection(sectionID)
    const urlSearch = useURLSearch()
    const q = urlSearch.get('q')
    const inDicts = urlSearch.get('in')
    const page = urlSearch.get('page')
    const dicts = useDictsInSection(sectionID)
    const searchState = useSearchState()

    const dispatch = useDispatch()
    useEffect(() => {
        dispatch(search(params, urlSearch))
        return () => { dispatch(searchReset()) }
    }, [sectionID, q, inDicts, page])

    const sections = useSections()

    if (!section) {
        return <NotFound />
    }

    const renderDictList = (): React.ReactNode => {
        let title = "Verbum - Анлайнавы Слоўнік Беларускай Мовы"
        if (sectionID !== 'default') {
            title = `${section.Name} - ${title}`
        }

        return <>
            <Helmet>
                <title>{title}</title>
                <meta name="description" content={title} />
                <meta property="og:title" content={title} />
                <meta property="og:description" content={title} />
                <meta name="robots" content="index, follow" />
            </Helmet>
            {section.Descr ? <p className='mx-1 mb-3'>{section.Descr}</p> : <></>}
            <DictsList dictionaries={dicts} />
        </>
    }

    const renderSearchResults = (): React.ReactNode => {
        let title = `${q} - Пошук`
        if (sectionID !== 'default') {
            title = `${title} - ${section.Name}`
        }

        return (
            <>
                <Helmet>
                    <title>{title}</title>
                    <meta name="description" content={title} />
                    <meta property="og:title" content={title} />
                    <meta property="og:description" content={title} />
                    <meta name="robots" content="noindex, nofollow" />
                </Helmet>
                {searchState.searchResult && searchState.searchResult.Articles.length > 0 && (
                    <> <div>{searchState.searchResult.Articles.map(
                        hit => (
                            <ArticleView
                                key={`${hit.DictionaryID}-${hit.ID}`}
                                a={hit}
                                showExternalButton={true}
                                showSource={true}
                            />
                        )
                    )}
                    </div>
                        <PaginationView
                            current={searchState.searchResult.Pagination.Current}
                            total={searchState.searchResult.Pagination.Total}
                            pageToURL={p => ({
                                search: urlSearch
                                    .clone()
                                    .set("page", p)
                                    .encode()
                            })}
                        />
                    </>
                )}
                {searchState.searchResult
                    && searchState.searchResult.Articles.length == 0
                    && renderNoSearchResults()
                }
            </>
        )
    }

    const renderNoSearchResults = (): React.ReactNode => <NoSearchResults q={q} suggestions={searchState.searchResult.TermSuggestions}
        calculateSuggestionURL={(s) => ({ search: urlSearch.clone().set('q', s).encode() })} />

    return (
        <div>
            <SearchControl
                inBound={dicts}
                urlQ={q}
                urlIn={inDicts}
                filterEnabled
                calculateSearchURL={
                    (q, inDicts) => urlSearch
                        .clone()
                        .set('q', q)
                        .set('in', inDicts)
                        .set('page', 1)
                        .encode()
                }
            />
            <ul className='nav nav-sections nav-underline mx-1 mb-1'>
                {sections.map((s, i) => <li className="nav-item" key={s.ID}>
                    <Link
                        className={`nav-link ${sectionID === s.ID ? 'active' : ''}`}
                        to={{
                            pathname: s.ID === 'default' ? '/' : `/s/${s.ID}`,
                            search: sectionID === s.ID ? '' : urlSearch.clone().reset('in').reset('page').encode(),
                        }}>{s.Name}</Link>
                </li>)}
            </ul>
            {!q ? renderDictList() : renderSearchResults()}
        </div>
    )
}
