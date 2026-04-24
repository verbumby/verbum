import { createSlice, type PayloadAction } from '@reduxjs/toolkit'
import { serverLoader } from '../../common/serverLoader'
import type { URLSearch } from '../../common/urlsearch'
import type { AppThunkAction } from '../../thunk'
import { type MatchParams, URLSearchDefaults } from './dict'

export type PrefaceState = string | null

const prefaceSlice = createSlice({
    name: 'preface',
    initialState: null as PrefaceState,
    reducers: {
        prefaceFetchKickOff: (state) => state,
        prefaceFetchSuccess: (_, action: PayloadAction<PrefaceState>) =>
            action.payload,
        prefaceFetchFailure: (state) => state,
        prefaceReset: () => null,
    },
})

const { prefaceFetchKickOff, prefaceFetchFailure } = prefaceSlice.actions
export const { prefaceFetchSuccess, prefaceReset } = prefaceSlice.actions
export const prefaceReducer = prefaceSlice.reducer

export const prefaceFetch = (
    params: MatchParams,
    urlSearch: URLSearch<typeof URLSearchDefaults>,
): AppThunkAction => {
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
            dispatch(prefaceFetchKickOff())
            dispatch(prefaceFetchSuccess(await verbumClient.getPreface(dictID)))
        } catch (err) {
            dispatch(prefaceFetchFailure())
            console.log('ERROR: ', err)
            throw err
        }
    }
}

export const prefaceFetchServer = serverLoader(URLSearchDefaults, prefaceFetch)
