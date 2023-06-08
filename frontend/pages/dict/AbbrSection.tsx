import * as React from 'react'
import { FC, useEffect } from 'react'
import { Helmet } from 'react-helmet'
import { useRouteMatch } from 'react-router-dom'
import { useDispatch } from '../../common'
import { useAbbr, useDict } from '../../store'
import { abbrFetch, abbrReset } from './abbr'
import { MatchParams, useURLSearch } from './dict'

export const AbbrSection: FC = ({}) => {
    const match = useRouteMatch<MatchParams>()
    const urlSearch = useURLSearch()

    const [dict] = useDict(match.params.dictID)
	const title = `Скарачэнні - ${dict.Title}`

    const abbr = useAbbr()
    const dispatch = useDispatch()

    useEffect(() => {
        dispatch(abbrFetch(match, urlSearch))
    }, [match.params.dictID])
    useEffect(() => () => { dispatch(abbrReset()) }, [])

    if (!abbr) {
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
				<h5>Скарачэнні</h5>
				<ul>
					{abbr.map(a => <li>{a.Keys.join(", ")} → <span className='text-secondary'>{a.Value}</span></li>)}
				</ul>
            </div>
		</>
	)
}
