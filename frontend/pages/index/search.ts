import { createSlice, type PayloadAction } from '@reduxjs/toolkit'
import { useURLSearch as useURLSearchCommon } from '../../common/hooks'
import type { SearchResult } from '../../common/search'
import { URLSearch } from '../../common/urlsearch'
import type { AppThunkAction } from '../../thunk'

export type MatchParams = {
    sectionID?: string
}

const URLSearchDefaults = {
    q: '',
    in: '',
    page: 1,
}

export const useURLSearch = () => useURLSearchCommon(URLSearchDefaults)

export type SearchState = {
    q: string
    searchResult: SearchResult
}

const searchSlice = createSlice({
    name: 'search',
    initialState: { q: '', searchResult: null } as SearchState,
    reducers: {
        searchKickOff: (state, action: PayloadAction<string>) => {
            state.q = action.payload
            state.searchResult = null
        },
        searchSuccess: (state, action: PayloadAction<SearchResult>) => {
            state.searchResult = action.payload
        },
        searchFailure: () => {},
        searchReset: () => ({ q: '', searchResult: null }),
    },
})

const { searchKickOff, searchSuccess, searchFailure } = searchSlice.actions
export const { searchReset } = searchSlice.actions
export const searchReducer = searchSlice.reducer

export const search = (
    params: Partial<MatchParams>,
    urlSearch: URLSearch<typeof URLSearchDefaults>,
): AppThunkAction => {
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
                const section = getState().sections.find(
                    (s) => s.ID === sectionID,
                )
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

export const searchServer = (
    params: Partial<MatchParams>,
    urlSearchParams: URLSearchParams,
): AppThunkAction =>
    search(params, new URLSearch(URLSearchDefaults, urlSearchParams))
