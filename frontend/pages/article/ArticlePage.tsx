import * as React from 'react'
import { Helmet } from 'react-helmet'
import { Redirect, useRouteMatch } from 'react-router-dom'
import { ArticleView, SearchControl, NotFound, useDispatch } from '../../common'
import { useArticle, useDict } from '../../store'
import { useURLSearch as useIndexURLSearch } from '../index/search'
import { articleFetch, articleReset, MatchParams } from './article'

export const ArticlePage: React.FC = () => {
    const match = useRouteMatch<MatchParams>()
    const [dict, dictIsAlias] = useDict(match.params.dictID)
    if (dictIsAlias) {
        return <Redirect to={{pathname: `/${dict.ID}/${match.params.articleID}` }} />
    }
    if (dict === null) {
        return <NotFound />
    }

    const a = useArticle()
    const indexURLSearch = useIndexURLSearch()

    const dispatch = useDispatch()
    React.useEffect(() => {
        if (a === undefined) {
            dispatch(articleFetch(match))
        }
        return () => { dispatch(articleReset()) }
    }, [match.path])

    if (a === undefined) {
        return <></>
    }

    if (a === null) {
        return <NotFound />
    }

    return (
        <>
            <Helmet>
                <title>{a.Title} - {dict.Title}</title>
                <meta name="description" content={`${a.Title} - ${dict.Title}`} />
                <meta property="og:title" content={`${a.Title} - ${dict.Title}`} />
                <meta property="og:description" content={`${a.Title} - ${dict.Title}`} />
                <meta name="robots" content="index, nofollow" />
            </Helmet>
            <div>
                <SearchControl
                    urlQ={a.Headword[0]}
                    urlIn=""
                    calculateSearchURL={
                        (q, inDicts) => indexURLSearch
                            .clone()
                            .set('q', q)
                            .set('in', inDicts)
                            .set('page', 1)
                            .encode()
                    }
                />
                <ArticleView a={a} showExternalButton={false} showSource={true} />
            </div>
        </>
    )
}
