import * as React from 'react'
import { Helmet } from 'react-helmet'
import { useDispatch } from 'react-redux'
import { useRouteMatch } from 'react-router-dom'
import { ArticleView, useURLSearch } from '../../common'
import { useArticle, useDicts } from '../../store'
import { articleFetch, articleReset, MatchParams } from './article'

export const ArticlePage: React.VFC = () => {
    const dicts = useDicts()
    const match = useRouteMatch<MatchParams>()
    const urlSearch = useURLSearch()
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
                <title>{a.Title} - {dicts.find(d => d.ID === a.DictionaryID).Title}</title>
                <meta name="robots" content="index, nofollow" />
            </Helmet>
            <div>
                <ArticleView a={a} />
            </div>
        </>
    )
}