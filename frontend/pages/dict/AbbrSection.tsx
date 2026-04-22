import * as React from 'react'
import { type FC, useEffect } from 'react'
import { Helmet } from 'react-helmet'
import { useParams } from 'react-router'
import { useDispatch } from '../../common/hooks'
import { useAbbr, useDictMust } from '../../store'
import { abbrFetch, abbrReset } from './abbr'
import { type MatchParams, useURLSearch } from './dict'

export const AbbrSection: FC = ({}) => {
    const params = useParams<MatchParams>()
    const urlSearch = useURLSearch()

    const [dict, _] = useDictMust(params.dictID)
    const title = `Скарачэнні - ${dict.Title}`

    const abbr = useAbbr()
    const dispatch = useDispatch()

    useEffect(() => {
        dispatch(abbrFetch(params, urlSearch))
    }, [params.dictID])
    useEffect(
        () => () => {
            dispatch(abbrReset())
        },
        [],
    )

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
            <div className="mx-1 mb-3">
                <h4>{dict.Title}</h4>
                <h5>Скарачэнні</h5>
                <ul>
                    {abbr.map((a) => (
                        <li>
                            {a.Keys.join(', ')} →{' '}
                            <span className="text-secondary">{a.Value}</span>
                        </li>
                    ))}
                </ul>
            </div>
        </>
    )
}
