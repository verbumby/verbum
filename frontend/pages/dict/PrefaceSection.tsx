import * as React from 'react'
import { FC, useEffect } from 'react'
import { Helmet } from 'react-helmet'
import { useDispatch } from '../../common'
import { usePreface, useDict } from '../../store'
import { prefaceFetch, prefaceReset } from './preface'
import { MatchParams, useURLSearch } from './dict'
import { useParams } from 'react-router'

export const PrefaceSection: FC = ({}) => {
    const params = useParams<MatchParams>()
    const urlSearch = useURLSearch()

    const [dict] = useDict(params.dictID)
	const title = `Прадмова - ${dict.Title}`

    const preface = usePreface()
    const dispatch = useDispatch()

    useEffect(() => {
        dispatch(prefaceFetch(params, urlSearch))
    }, [params.dictID])
    useEffect(() => () => { dispatch(prefaceReset()) }, [])

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
            <div className='mx-1 mb-3'>
				<div dangerouslySetInnerHTML={{ __html: preface }} />
            </div>
		</>
	)
}
