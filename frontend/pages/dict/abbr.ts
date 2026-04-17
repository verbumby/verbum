import { createSlice, type PayloadAction } from '@reduxjs/toolkit'
import type { Abbrevs } from '../../common/abbrevs'
import { URLSearch } from '../../common/urlsearch'
import type { AppThunkAction } from '../../thunk'
import { type MatchParams, URLSearchDefaults } from './dict'

export type AbbrState = Abbrevs

const abbrSlice = createSlice({
    name: 'abbr',
    initialState: null as AbbrState,
    reducers: {
        abbrFetchKickOff: (state) => state,
        abbrFetchSuccess: (_, action: PayloadAction<Abbrevs>) => action.payload,
        abbrFetchFailure: (state) => state,
        abbrReset: () => null,
    },
})

const { abbrFetchKickOff, abbrFetchFailure } = abbrSlice.actions
export const { abbrFetchSuccess, abbrReset } = abbrSlice.actions
export const abbrReducer = abbrSlice.reducer

export const abbrFetch = (
    params: Partial<MatchParams>,
    urlSearch: URLSearch<typeof URLSearchDefaults>,
): AppThunkAction => {
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

export const abbrFetchServer = (
    params: Partial<MatchParams>,
    urlSearchParams: URLSearchParams,
): AppThunkAction =>
    abbrFetch(params, new URLSearch(URLSearchDefaults, urlSearchParams))
