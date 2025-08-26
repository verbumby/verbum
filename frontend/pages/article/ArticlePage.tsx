import * as React from 'react'
import { Helmet } from 'react-helmet'
import { Redirect, useRouteMatch } from 'react-router-dom'
import { ArticleView, SearchControl, NotFound, useDispatch } from '../../common'
import { useArticle, useDict } from '../../store'
import { useURLSearch as useDictURLSearch } from '../dict/dict'
import { articleFetch, articleReset, MatchParams } from './article'

export const ArticlePage: React.FC = () => {
    const match = useRouteMatch<MatchParams>()
    const [dict, dictIsAlias] = useDict(match.params.dictID)
    if (dictIsAlias) {
        return <Redirect to={{ pathname: `/${dict.ID}/${match.params.articleID}` }} />
    }
    if (dict === null) {
        return <NotFound />
    }

    const a = useArticle()
    const dictURLSearch = useDictURLSearch()

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

    const title = `${a.Title} - ${dict.Title}`
    return (
        <>
            <Helmet>
                <title>{title}</title>
                <meta name="description" content={title} />
                <meta property="og:title" content={title} />
                <meta property="og:description" content={title} />
                <meta name="robots" content="index, nofollow" />
            </Helmet>
            <div>
                <SearchControl
                    inBound={[dict]}
                    urlQ={a.Headword[0]}
                    urlIn={dict.ID}
                    filterEnabled={false}
                    calculateSearchURL={
                        (q, inDicts) => ({
                            pathname: `/${dict.ID}`,
                            search: dictURLSearch.clone().set('q', q).encode(),
                        })
                    }
                />
                <ArticleView a={a} showExternalButton={false} showSource={true} />
            </div>
        </>
    )
}
