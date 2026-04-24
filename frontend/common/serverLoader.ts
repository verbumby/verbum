import type { Params } from 'react-router'
import type { AppThunkAction } from '../thunk'
import type { URLSearchEntries } from './urlsearch'
import { URLSearch } from './urlsearch'

export function serverLoader<P, D extends URLSearchEntries>(
    defaults: D,
    action: (params: P, urlSearch: URLSearch<D>) => AppThunkAction,
): (
    params: Params<string>,
    urlSearchParams: URLSearchParams,
) => AppThunkAction {
    return (params, urlSearchParams) =>
        action(params as P, new URLSearch(defaults, urlSearchParams))
}
