import { combineReducers } from 'redux'
import { useSelector as useSelectorParent } from 'react-redux'
import { SearchActions, SearchState, searchReducer } from './pages/index/index'
import { Dict, Article, DictsActions, DictsState, dictsReducer } from './common'
import { ThunkAction } from '@reduxjs/toolkit'
import { ArticleState, ArticleActions, articleReducer } from './pages/article'
import { DictArticlesActions, LetterFilterActions, letterFilterReducer, LetterFilterState, DictArticlesState, dictArticlesReducer } from './pages/dict'
import { loadingBarReducer } from 'react-redux-loading-bar'

type AllActions = DictsActions | SearchActions | ArticleActions | LetterFilterActions | DictArticlesActions

export type RootState = {
    dicts: DictsState
    search: SearchState
    article: ArticleState
    letterFilter: LetterFilterState
    dictArticles: DictArticlesState
    loadingBar: any
}

export const rootReducer = combineReducers<RootState, AllActions>({
    dicts: dictsReducer,
    search: searchReducer,
    article: articleReducer,
    letterFilter: letterFilterReducer,
    dictArticles: dictArticlesReducer,
    loadingBar: loadingBarReducer,
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

export function useDict(id: string): [Dict, boolean] {
    let d = useDicts().find(d => d.ID === id)
    if (d) {
        return [d, false]
    }
    d = useDicts().find(d => d.Aliases && d.Aliases.includes(id))
    if (d) {
        return [d, true]
    }
    return [null, false]
}

export function useSearchState(): SearchState {
    return useSelector<SearchState>(state => state.search)
}

export function useArticleState(): ArticleState {
    return useSelector<ArticleState>(state => state.article)
}

export function useArticle(): Article {
    return useArticleState().a
}

export function useLetterFilter(): LetterFilterState {
    return useSelector<LetterFilterState>(state => state.letterFilter)
}

export function useDictArticles(): DictArticlesState {
    return useSelector<DictArticlesState>(state => state.dictArticles)
}

export type AppThunkAction<ReturnType = void> = ThunkAction<Promise<ReturnType>, RootState, unknown, AllActions>
