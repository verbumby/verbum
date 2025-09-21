import { Abbrevs, URLSearch } from "../../common"
import { AppThunkAction } from "../../store"
import { MatchParams, URLSearchDefaults } from './dict'

export type AbbrState = Abbrevs

const ABBR_FETCH_KICKOFF = 'ABBR/FETCH/KICKOFF'
type AbbrFetchKickOffAction = {
    type: typeof ABBR_FETCH_KICKOFF
    dictID: string
}
function abbrFetchKickOff(dictID: string): AbbrFetchKickOffAction {
    return {
        type: ABBR_FETCH_KICKOFF,
        dictID,
    }
}

const ABBR_FETCH_SUCCESS = 'ABBR/FETCH/SUCCESS'
type AbbrFetchSuccessAction = {
    type: typeof ABBR_FETCH_SUCCESS
    abbr: Abbrevs
}
function abbrFetchSuccess(abbr: Abbrevs): AbbrFetchSuccessAction {
    return {
        type: ABBR_FETCH_SUCCESS,
        abbr,
    }
}

const ABBR_FETCH_FAILURE = 'ABBR/FETCH/FAILURE'
type AbbrFetchFailureAction = {
    type: typeof ABBR_FETCH_FAILURE
}
function abbrFetchFailure(): AbbrFetchFailureAction {
    return { type: ABBR_FETCH_FAILURE }
}

const ABBR_RESET = 'ABBR/RESET'
type AbbrResetAction = {
    type: typeof ABBR_RESET
}
export function abbrReset(): AbbrResetAction {
    return { type: ABBR_RESET }
}

export type AbbrActions = AbbrFetchKickOffAction | AbbrFetchSuccessAction | AbbrFetchFailureAction | AbbrResetAction

export function abbrReducer(state: AbbrState = null, a: AbbrActions): AbbrState {
    switch (a.type) {
        case ABBR_FETCH_KICKOFF:
            return state
        case ABBR_FETCH_SUCCESS:
            return a.abbr
        case ABBR_RESET:
            return null
        default:
            return state
    }
}

export const abbrFetch = (params: Partial<MatchParams>, urlSearch: URLSearch<typeof URLSearchDefaults>): AppThunkAction => {
    return async (dispatch, getState): Promise<void> => {
        try {
            if (urlSearch.get('section') !== 'abbr') {
                return
            }

            const state = getState()
            if (state.abbr !== null) {
                return
            }

            const { dictID } = params
            dispatch(abbrFetchKickOff(dictID))
            dispatch(abbrFetchSuccess(await verbumClient.getAbbr(dictID)))
        } catch (err) {
            dispatch(abbrFetchFailure())
            console.log('ERROR: ', err)
            throw err
        }
    }
}

export const abbrFetchServer = (params: Partial<MatchParams>, urlSearchParams: URLSearchParams): AppThunkAction =>
    abbrFetch(params, new URLSearch(URLSearchDefaults, urlSearchParams))
