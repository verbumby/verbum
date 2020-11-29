import { combineReducers } from 'redux'
import { useSelector as useSelectorParent } from 'react-redux'
import { SearchActions, SearchState, searchReducer } from './pages/index/index'
import { DictsActions, DictsState, dictsReducer } from './common'

type AllActions = DictsActions | SearchActions

export type RootState = {
    dicts: DictsState
    search: SearchState
}

export const rootReducer = combineReducers<RootState, AllActions>({
    dicts: dictsReducer,
    search: searchReducer,
})

export function useSelector<TSelected = unknown>(
    selector: (state: RootState) => TSelected,
    equalityFn?: (left: TSelected, right: TSelected) => boolean
): TSelected {
    return useSelectorParent<RootState,TSelected>(selector, equalityFn)
}

export function useDicts(): DictsState {
    return useSelector<DictsState>(state => state.dicts)
}

export function useSearchState(): SearchState {
    return useSelector<SearchState>(state => state.search)
}
