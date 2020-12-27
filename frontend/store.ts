import { combineReducers } from 'redux'
import { useSelector as useSelectorParent } from 'react-redux'
import { SearchActions, SearchState, searchReducer } from './pages/index/index'
import { Dict, DictsActions, DictsState, dictsReducer } from './common'
import { ThunkAction } from '@reduxjs/toolkit'
import { ArticleState, ArticleActions, articleReducer } from './pages/article'
import { DictArticlesActions, LetterFilterActions, letterFilterReducer, LetterFilterState, DictArticlesState, dictArticlesReducer } from './pages/dict'

type AllActions = DictsActions | SearchActions | ArticleActions | LetterFilterActions | DictArticlesActions

export type RootState = {
    dicts: DictsState
    search: SearchState
    article: ArticleState
    letterFilter: LetterFilterState
    dictArticles: DictArticlesState
}

export const rootReducer = combineReducers<RootState, AllActions>({
    dicts: dictsReducer,
    search: searchReducer,
    article: articleReducer,
    letterFilter: letterFilterReducer,
    dictArticles: dictArticlesReducer,
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

export function useDict(id: string): Dict {
    return useDicts().find(d => d.ID === id)
}

export function useSearchState(): SearchState {
    return useSelector<SearchState>(state => state.search)
}

export function useArticle(): ArticleState {
    return useSelector<ArticleState>(state => state.article)
}

export function useLetterFilter(): LetterFilterState {
    return useSelector<LetterFilterState>(state => state.letterFilter)
}

export function useDictArticles(): DictArticlesState {
    return useSelector<DictArticlesState>(state => state.dictArticles)
}

export type AppThunkAction<ReturnType = void> = ThunkAction<Promise<ReturnType>, RootState, unknown, AllActions>
