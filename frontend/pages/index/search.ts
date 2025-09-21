import { SearchResult, URLSearch, useURLSearch as useURLSearchCommon } from '../../common'
import { AppThunkAction } from '../../store'

export type MatchParams = {
    sectionID?: string,
}

const URLSearchDefaults = {
    q: '',
    in: '',
    page: 1
}

export const useURLSearch = () => useURLSearchCommon(URLSearchDefaults)

export type SearchState = {
    q: string
    searchResult: SearchResult
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
    searchResult: SearchResult
}
function searchSuccess(searchResult: SearchResult): SearchSuccessAction {
    return {
        type: SEARCH_SUCCESS,
        searchResult,
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

export function searchReducer(state: SearchState = {q: '', searchResult: null}, a: SearchActions): SearchState {
    switch (a.type) {
        case SEARCH_KICKOFF:
            return {
                q: a.q,
                searchResult: null,
            }
        case SEARCH_SUCCESS:
            return {
                ...state,
                searchResult: a.searchResult,
            }
        case SEARCH_RESET:
            return {
                q: '',
                searchResult: null,
            }
        default:
            return state
    }
}

export const search = (params: Partial<MatchParams>, urlSearch: URLSearch<typeof URLSearchDefaults>): AppThunkAction => {
    return async (dispatch, getState): Promise<void> => {
        try {
            const q = urlSearch.get('q')
            if (!q) {
                return
            }
            if (q === getState().search.q) {
                return
            }

            let inDicts = urlSearch.get('in')
            if (!inDicts) {
                const sectionID = params.sectionID || 'default'
                const section = getState().sections.find(s => s.ID === sectionID)
                if (!section) {
                    return
                }
                inDicts = section.DictIDs.join(',')
            }
            const page = urlSearch.get('page')

            dispatch(searchKickOff(q))
            dispatch(searchSuccess(await verbumClient.search(q, inDicts, page)))
        } catch (err) {
            dispatch(searchFailure())
            console.log('ERROR: ', err)
            throw err
        }
    }
}

export const searchServer = (params: Partial<MatchParams>, urlSearchParams: URLSearchParams): AppThunkAction =>
    search(params, new URLSearch(URLSearchDefaults, urlSearchParams))
