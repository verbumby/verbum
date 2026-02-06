import * as React from 'react'
import { Helmet } from 'react-helmet'
import { useLocation, useParams } from 'react-router'
import { ArticleView, SearchControl, NotFound, useDispatch, Redirect } from '../../common'
import { useArticle, useDict } from '../../store'
import { useURLSearch as useDictURLSearch } from '../dict/dict'
import { articleFetch, articleReset, MatchParams, useURLSearch } from './article'

export const ArticlePage: React.FC = () => {
    const params = useParams<MatchParams>()
    const location = useLocation()
    const [dict, dictIsAlias] = useDict(params.dictID)
    const urlSearch = useURLSearch()
    const ogPreview = !!urlSearch.get('og_preview')

    if (dictIsAlias) {
        return <Redirect to={{ pathname: `/${dict.ID}/${params.articleID}` }} />
    }

    if (dict === null) {
        return <NotFound />
    }

    const a = useArticle()
    const dictURLSearch = useDictURLSearch()

    const dispatch = useDispatch()
    React.useEffect(() => {
        if (a === undefined) {
            dispatch(articleFetch(params))
        }
        return () => { dispatch(articleReset()) }
    }, [location])

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
                <meta property="og:image" content={`/api/dictionaries/${encodeURIComponent(dict.ID)}/articles/${encodeURIComponent(params.articleID)}/preview`} />
                <meta name="robots" content="index, nofollow" />
            </Helmet>
            <div>
                {!ogPreview &&
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
                }
                <ArticleView a={a} showExternalButton={false} showCopyURLButton={!ogPreview} showSource={true} />
            </div>
        </>
    )
}
