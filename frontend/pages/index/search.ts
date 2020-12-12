import { match } from 'react-router-dom'

import { Article } from '../../common/article'
import { AppThunkAction } from '../../store'

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

const SEARCH_SET_HITS = 'SEARCH/SET_HITS'
type SearchSetHitsAction = {
    type: typeof SEARCH_SET_HITS
    hits: Article[]
}
function searchSetHits(hits: Article[]): SearchSetHitsAction {
    return {
        type: SEARCH_SET_HITS,
        hits,
    }
}

const SEARCH_RESET = 'SEARCH/RESET'
type SearchResetAction = {
    type: typeof SEARCH_RESET,
}
export function searchReset(): SearchResetAction {
    return {type: SEARCH_RESET}
}

export type SearchActions = SearchKickOffAction | SearchSetHitsAction | SearchResetAction

export function searchReducer(state: SearchState = {q: '', hits: []}, a: SearchActions): SearchState {
    switch (a.type) {
        case SEARCH_KICKOFF:
            return {
                q: a.q,
                hits: [],
            }
        case SEARCH_SET_HITS:
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

export const search = (match: match, urlSearch: URLSearchParams): AppThunkAction => {
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
            dispatch(searchSetHits(await verbumClient.search(q)))
        } catch (err) {
            console.log('ERROR: ', err)
            throw err
        }
    }
}
