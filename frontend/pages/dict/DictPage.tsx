import * as React from 'react'
import { Redirect, useRouteMatch } from 'react-router-dom'
import { NotFound } from '../../common'
import { useDict } from '../../store'
import { AbbrSection } from './AbbrSection'
import { DefaultSection } from './DefaultSection'
import { MatchParams, useURLSearch } from './dict'

export const DictPage: React.FC = ({ }) => {
    const match = useRouteMatch<MatchParams>()
    const urlSearch = useURLSearch()

    const [dict, dictIsAlias] = useDict(match.params.dictID)
    if (dictIsAlias) {
        return <Redirect to={{pathname: `/${dict.ID}`, search: urlSearch.encode() }} />
    }
    if (dict === null) {
        return <NotFound />
    }

    if (urlSearch.get('section') == 'abbr') {
        return <AbbrSection />
    }

    return <DefaultSection />
}
