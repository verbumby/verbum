import { createSlice, type PayloadAction } from '@reduxjs/toolkit'
import type { LetterFilter } from '../../common/letterfilter'
import { serverLoader } from '../../common/serverLoader'
import type { URLSearch } from '../../common/urlsearch'
import type { AppThunkAction } from '../../thunk'
import { type MatchParams, URLSearchDefaults } from './dict'

export type LetterFilterState = LetterFilter | null

const letterFilterSlice = createSlice({
    name: 'letterFilter',
    initialState: null as LetterFilterState,
    reducers: {
        letterFilterFetchKickOff: (state) => state,
        letterFilterFetchSuccess: (_, action: PayloadAction<LetterFilter>) =>
            action.payload,
        letterFilterFetchFailure: (state) => state,
        letterFilterReset: () => null,
    },
})

const { letterFilterFetchKickOff, letterFilterFetchFailure } =
    letterFilterSlice.actions
export const { letterFilterFetchSuccess, letterFilterReset } =
    letterFilterSlice.actions
export const letterFilterReducer = letterFilterSlice.reducer

export const letterFilterFetch = (
    params: MatchParams,
    urlSearch: URLSearch<typeof URLSearchDefaults>,
): AppThunkAction => {
    return async (dispatch, getState): Promise<void> => {
        try {
            if (urlSearch.get('section') !== '') {
                return
            }

            const { dictID } = params
            const prefix = urlSearch.get('prefix')
            const state = getState()
            if (
                state.letterFilter &&
                state.letterFilter.DictID === dictID &&
                state.letterFilter.Prefix === prefix
            ) {
                return
            }
            dispatch(letterFilterFetchKickOff())
            dispatch(
                letterFilterFetchSuccess(
                    await verbumClient.getLetterFilter(dictID, prefix),
                ),
            )
        } catch (err) {
            dispatch(letterFilterFetchFailure())
            console.log('ERROR: ', err)
            throw err
        }
    }
}

export const letterFilterFetchServer = serverLoader(
    URLSearchDefaults,
    letterFilterFetch,
)
