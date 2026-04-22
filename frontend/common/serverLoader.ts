import type { AppThunkAction } from '../thunk'
import type { Params } from 'react-router'
import { URLSearch, type URLSearchEntries } from './urlsearch'

export function serverLoader<P, D extends URLSearchEntries>(
    defaults: D,
    action: (params: P, urlSearch: URLSearch<D>) => AppThunkAction,
): (params: Params<string>, urlSearchParams: URLSearchParams) => AppThunkAction {
    return (params, urlSearchParams) => action(params as P, new URLSearch(defaults, urlSearchParams))
}
