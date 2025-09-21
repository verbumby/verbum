import { URLSearch } from "../../common"
import { AppThunkAction } from "../../store"
import { MatchParams, URLSearchDefaults } from './dict'

export type PrefaceState = string | null

const PREFACE_FETCH_KICKOFF = 'PREFACE/FETCH/KICKOFF'
type PrefaceFetchKickOffAction = {
    type: typeof PREFACE_FETCH_KICKOFF
    dictID: string
}
function prefaceFetchKickOff(dictID: string): PrefaceFetchKickOffAction {
    return {
        type: PREFACE_FETCH_KICKOFF,
        dictID,
    }
}

const PREFACE_FETCH_SUCCESS = 'PREFACE/FETCH/SUCCESS'
type PrefaceFetchSuccessAction = {
    type: typeof PREFACE_FETCH_SUCCESS
    preface: PrefaceState
}
function prefaceFetchSuccess(preface: PrefaceState): PrefaceFetchSuccessAction {
    return {
        type: PREFACE_FETCH_SUCCESS,
        preface,
    }
}

const PREFACE_FETCH_FAILURE = 'PREFACE/FETCH/FAILURE'
type PrefaceFetchFailureAction = {
    type: typeof PREFACE_FETCH_FAILURE
}
function prefaceFetchFailure(): PrefaceFetchFailureAction {
    return { type: PREFACE_FETCH_FAILURE }
}

const PREFACE_RESET = 'PREFACE/RESET'
type PrefaceResetAction = {
    type: typeof PREFACE_RESET
}
export function prefaceReset(): PrefaceResetAction {
    return { type: PREFACE_RESET }
}

export type PrefaceActions = PrefaceFetchKickOffAction | PrefaceFetchSuccessAction | PrefaceFetchFailureAction | PrefaceResetAction

export function prefaceReducer(state: PrefaceState = null, a: PrefaceActions): PrefaceState {
    switch (a.type) {
        case PREFACE_FETCH_KICKOFF:
            return state
        case PREFACE_FETCH_SUCCESS:
            return a.preface
        case PREFACE_RESET:
            return null
        default:
            return state
    }
}

export const prefaceFetch = (params: Partial<MatchParams>, urlSearch: URLSearch<typeof URLSearchDefaults>): AppThunkAction => {
    return async (dispatch, getState): Promise<void> => {
        try {
            if (urlSearch.get('section') !== 'preface') {
                return
            }

            const state = getState()
            if (state.preface !== null) {
                return
            }

            const { dictID } = params
            dispatch(prefaceFetchKickOff(dictID))
            dispatch(prefaceFetchSuccess(await verbumClient.getPreface(dictID)))
        } catch (err) {
            dispatch(prefaceFetchFailure())
            console.log('ERROR: ', err)
            throw err
        }
    }
}

export const prefaceFetchServer = (params: Partial<MatchParams>, urlSearchParams: URLSearchParams): AppThunkAction =>
    prefaceFetch(params, new URLSearch(URLSearchDefaults, urlSearchParams))
