import * as React from 'react'
import { useEffect } from 'react'
import { NoSearchResults, NotFound, useDispatch } from '../../common'
import { Link, useRouteMatch } from 'react-router-dom'
import { Helmet } from "react-helmet"

import { useDictsInSection, useSearchState, useSection, useSections } from '../../store'
import { search, searchReset, useURLSearch } from './search'
import { DictsList } from './DictsList'
import { ArticleView, PaginationView, SearchControl } from '../../common'

export const IndexPage: React.FC = () => {
    const match = useRouteMatch<{ sectionID?: string }>()
    const sectionID = match.params.sectionID || 'default'
    const section = useSection(sectionID)
    const urlSearch = useURLSearch()
    const q = urlSearch.get('q')
    const inDicts = urlSearch.get('in')
    const page = urlSearch.get('page')
    const dicts = useDictsInSection(sectionID)
    const searchState = useSearchState()

    const dispatch = useDispatch()
    useEffect(() => {
        dispatch(search(match, urlSearch))
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
            <p />
            <h4 className='mx-1 mb-3'>{section.Name}</h4>
            <p className='mx-1 mb-3'>{section.Descr}</p>
            <DictsList dictionaries={dicts} />
            <p className='mx-1 mb-3'>Іншыя раздзелы: {
                sections.filter(s => s.ID !== sectionID).map((s, i) => <React.Fragment key={s.ID}>
                    {i == 0 ? '' : ', '}
                    <Link to={s.ID === 'default' ? '/' : `/s/${s.ID}`}>{s.Name}</Link>
                </React.Fragment>)
            }</p>
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
            {!q ? renderDictList() : renderSearchResults()}
        </div>
    )
}
