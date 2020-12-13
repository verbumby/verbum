import * as React from 'react'
import { Helmet } from 'react-helmet'
import { useDispatch } from 'react-redux'
import { useRouteMatch } from 'react-router-dom'
import { ArticleView, useURLSearch, SearchControl } from '../../common'
import { useArticle, useDict } from '../../store'
import { articleFetch, articleReset, MatchParams } from './article'

export const ArticlePage: React.VFC = () => {
    const match = useRouteMatch<MatchParams>()
    const urlSearch = useURLSearch()
    const dict = useDict(match.params.dictID)
    const a = useArticle()

    const dispatch = useDispatch()
    React.useEffect(() => {
        if (!a) {
            dispatch(articleFetch(match, urlSearch))
        }
        return () => dispatch(articleReset())
    }, [match.path])

    if (!a) {
        return <></>
    }

    return (
        <>
            <Helmet>
                <title>{a.Title} - {dict.Title}</title>
                <meta name="robots" content="index, nofollow" />
            </Helmet>
            <div>
                <SearchControl urlQ={a.Headword[0]} />
                <ArticleView a={a} />
            </div>
        </>
    )
}
