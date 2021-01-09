import { match } from "react-router-dom"
import { LetterFilter, URLSearch } from "../../common"
import { AppThunkAction } from "../../store"
import { MatchParams, URLSearchDefaults } from './dict'

export type LetterFilterState = LetterFilter

const LETTER_FILTER_FETCH_KICKOFF = 'LETTER_FILTER/FETCH/KICKOFF'
type LetterFilterFetchKickOffAction = {
    type: typeof LETTER_FILTER_FETCH_KICKOFF
    dictID: string
    prefix: string
}
function letterFilterFetchKickOff(dictID: string, prefix: string): LetterFilterFetchKickOffAction {
    return {
        type: LETTER_FILTER_FETCH_KICKOFF,
        dictID,
        prefix,
    }
}

const LETTER_FILTER_FETCH_SUCCESS = 'LETTER_FILTER/FETCH/SUCCESS'
type LetterFilterFetchSuccessAction = {
    type: typeof LETTER_FILTER_FETCH_SUCCESS
    letterFilter: LetterFilter
}
function letterFilterFetchSuccess(letterFilter: LetterFilter): LetterFilterFetchSuccessAction {
    return {
        type: LETTER_FILTER_FETCH_SUCCESS,
        letterFilter,
    }
}

const LETTER_FILTER_FETCH_FAILURE = 'LETTER_FILTER/FETCH/FAILURE'
type LetterFilterFetchFailureAction = {
    type: typeof LETTER_FILTER_FETCH_FAILURE
}
function letterFilterFetchFailure(): LetterFilterFetchFailureAction {
    return { type: LETTER_FILTER_FETCH_FAILURE }
}

const LETTER_FILTER_RESET = 'LETTER_FILTER/RESET'
type LetterFilterResetAction = {
    type: typeof LETTER_FILTER_RESET
}
export function letterFilterReset(): LetterFilterResetAction {
    return { type: LETTER_FILTER_RESET }
}

export type LetterFilterActions = LetterFilterFetchKickOffAction | LetterFilterFetchSuccessAction | LetterFilterFetchFailureAction | LetterFilterResetAction

export function letterFilterReducer(state: LetterFilterState = null, a: LetterFilterActions): LetterFilterState {
    switch (a.type) {
        case LETTER_FILTER_FETCH_KICKOFF:
            return state
        case LETTER_FILTER_FETCH_SUCCESS:
            return a.letterFilter
        case LETTER_FILTER_RESET:
            return null
        default:
            return state
    }
}

export const letterFilterFetch = (match: match<MatchParams>, urlSearch: URLSearch<typeof URLSearchDefaults>): AppThunkAction => {
    return async (dispatch, getState): Promise<void> => {
        try {
            const { dictID } = match.params
            const prefix = urlSearch.get('prefix')
            dispatch(letterFilterFetchKickOff(dictID, prefix))
            dispatch(letterFilterFetchSuccess(await verbumClient.getLetterFilter(dictID, prefix)))
        } catch (err) {
            dispatch(letterFilterFetchFailure())
            console.log('ERROR: ', err)
            throw err
        }
    }
}

export const letterFilterFetchServer = (match: match<MatchParams>, urlSearchParams: URLSearchParams): AppThunkAction =>
    letterFilterFetch(match, new URLSearch(URLSearchDefaults, urlSearchParams))
