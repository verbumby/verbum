import * as React from 'react'
import { type FC, useEffect } from 'react'
import { Helmet } from 'react-helmet'
import { useParams } from 'react-router'
import { useDispatch } from '../../common/hooks'
import { useDictMust, usePreface } from '../../store'
import { type MatchParams, useURLSearch } from './dict'
import { prefaceFetch, prefaceReset } from './preface'

export const PrefaceSection: FC = ({}) => {
    const params = useParams() as MatchParams
    const urlSearch = useURLSearch()

    const [dict, _] = useDictMust(params.dictID)
    const title = `Прадмова - ${dict.Title}`

    const preface = usePreface()
    const dispatch = useDispatch()

    useEffect(() => {
        dispatch(prefaceFetch(params, urlSearch))
    }, [params.dictID])
    useEffect(
        () => () => {
            dispatch(prefaceReset())
        },
        [],
    )

    if (!preface) {
        return <></>
    }

    return (
        <>
            <Helmet>
                <title>{title}</title>
                <meta name="description" content={title} />
                <meta property="og:title" content={title} />
                <meta property="og:description" content={title} />
                <meta name="robots" content="index, follow" />
            </Helmet>
            <div className="mx-1 mb-3">
                <div dangerouslySetInnerHTML={{ __html: preface }} />
            </div>
        </>
    )
}
