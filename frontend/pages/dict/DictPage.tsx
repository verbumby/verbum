import * as React from 'react'
import { useRouteMatch } from 'react-router-dom'
import { useDict } from '../../store'
import { MatchParams } from './dict'

export const DictPage: React.VFC = ({}) => {
    const match = useRouteMatch<MatchParams>()
    const dict = useDict(match.params.dictID)
    return <div>Dict Page - {dict.Title}</div>
}
