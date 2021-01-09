import { match } from 'react-router-dom'
import { URLSearchParams } from 'url'

import { URLSearch, useURLSearch as useURLSearchCommon } from '../../common'
import { Article } from '../../common/article'
import { AppThunkAction } from '../../store'

const URLSearchDefaults = {
    q: ''
}

export const useURLSearch = () => useURLSearchCommon(URLSearchDefaults)

export type SearchState = {
    q: string
    hits: Article[]
}

const SEARCH_KICKOFF = 'SEARCH/KICKOFF'
type SearchKickOffAction = {
    type: typeof SEARCH_KICKOFF,
    q: string
}
function searchKickOff(q: string): SearchKickOffAction {
    return {
        type: SEARCH_KICKOFF,
        q,
    }
}

const SEARCH_SUCCESS = 'SEARCH/SUCCESS'
type SearchSuccessAction = {
    type: typeof SEARCH_SUCCESS
    hits: Article[]
}
function searchSuccess(hits: Article[]): SearchSuccessAction {
    return {
        type: SEARCH_SUCCESS,
        hits,
    }
}

const SEARCH_FAILURE = 'SEARCH/FAILURE'
type SearchFailureAction = {
    type: typeof SEARCH_FAILURE
}
function searchFailure(): SearchFailureAction {
    return { type: SEARCH_FAILURE }
}

const SEARCH_RESET = 'SEARCH/RESET'
type SearchResetAction = {
    type: typeof SEARCH_RESET,
}
export function searchReset(): SearchResetAction {
    return {type: SEARCH_RESET}
}

export type SearchActions = SearchKickOffAction | SearchSuccessAction | SearchFailureAction | SearchResetAction

export function searchReducer(state: SearchState = {q: '', hits: []}, a: SearchActions): SearchState {
    switch (a.type) {
        case SEARCH_KICKOFF:
            return {
                q: a.q,
                hits: [],
            }
        case SEARCH_SUCCESS:
            return {
                ...state,
                hits: a.hits,
            }
        case SEARCH_RESET:
            return {
                q: '',
                hits: [],
            }
        default:
            return state
    }
}

export const search = (match: match, urlSearch: URLSearch<typeof URLSearchDefaults>): AppThunkAction => {
    return async (dispatch, getState): Promise<void> => {
        try {
            const q = urlSearch.get('q')
            if (!q) {
                return
            }
            if (q === getState().search.q) {
                return
            }

            dispatch(searchKickOff(q))
            dispatch(searchSuccess(await verbumClient.search(q)))
        } catch (err) {
            dispatch(searchFailure())
            console.log('ERROR: ', err)
            throw err
        }
    }
}

export const searchServer = (match: match, urlSearchParams: URLSearchParams): AppThunkAction =>
    search(match, new URLSearch(URLSearchDefaults, urlSearchParams))
