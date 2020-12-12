import { combineReducers } from 'redux'
import { useSelector as useSelectorParent } from 'react-redux'
import { SearchActions, SearchState, searchReducer } from './pages/index/index'
import { DictsActions, DictsState, dictsReducer, Article } from './common'
import { ThunkAction } from '@reduxjs/toolkit'
import { ArticleActions, articleReducer } from './pages/article'

type AllActions = DictsActions | SearchActions | ArticleActions

export type RootState = {
    dicts: DictsState
    search: SearchState
    article: Article
}

export const rootReducer = combineReducers<RootState, AllActions>({
    dicts: dictsReducer,
    search: searchReducer,
    article: articleReducer,
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

export function useArticle(): Article {
    return useSelector<Article>(state => state.article)
}

export type AppThunkAction<ReturnType = void> = ThunkAction<Promise<ReturnType>, RootState, unknown, AllActions>
