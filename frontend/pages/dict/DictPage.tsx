import type * as React from 'react'
import { useParams } from 'react-router'
import { NotFound } from '../../common/NotFound'
import { Redirect } from '../../common/Redirect'
import { useDict } from '../../store'
import { AbbrSection } from './AbbrSection'
import { DefaultSection } from './DefaultSection'
import { type MatchParams, useURLSearch } from './dict'
import { PrefaceSection } from './PrefaceSection'

export const DictPage: React.FC = ({}) => {
    const params = useParams<MatchParams>()
    const urlSearch = useURLSearch()

    const [dict, dictIsAlias] = useDict(params.dictID)
    if (dictIsAlias) {
        return (
            <Redirect
                to={{ pathname: `/${dict.ID}`, search: urlSearch.encode() }}
            />
        )
    }
    if (dict === null) {
        return <NotFound />
    }

    if (urlSearch.get('section') == 'abbr') {
        return <AbbrSection />
    }

    if (urlSearch.get('section') == 'preface') {
        return <PrefaceSection />
    }

    return <DefaultSection />
}
