import * as React from 'react'
import { DictionariesListState, useSelector } from '../reducers'

export const IndexPage: React.VFC = () => {
    const dicts = useSelector<DictionariesListState>(state => state.dictionaries)

    return <div>Index {JSON.stringify(dicts)} {dicts.length} Page</div>
}
