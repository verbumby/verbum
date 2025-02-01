import * as React from 'react'
import { FC, useEffect } from 'react'
import { Helmet } from 'react-helmet'
import { useRouteMatch } from 'react-router-dom'
import { useDispatch } from '../../common'
import { usePreface, useDict } from '../../store'
import { prefaceFetch, prefaceReset } from './preface'
import { MatchParams, useURLSearch } from './dict'

export const PrefaceSection: FC = ({}) => {
    const match = useRouteMatch<MatchParams>()
    const urlSearch = useURLSearch()

    const [dict] = useDict(match.params.dictID)
	const title = `Прадмова - ${dict.Title}`

    const preface = usePreface()
    const dispatch = useDispatch()

    useEffect(() => {
        dispatch(prefaceFetch(match, urlSearch))
    }, [match.params.dictID])
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
                <h4>{dict.Title}</h4>
				<h5>Прадмова</h5>
				<div dangerouslySetInnerHTML={{ __html: preface }} />
            </div>
		</>
	)
}
