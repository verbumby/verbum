import { combineReducers } from 'redux'
import { useSelector as useSelectorParent } from 'react-redux'
import { SearchActions, SearchState, searchReducer } from './pages/index/search'
import { Dict } from './common/dict'
import { Article } from './common/article'
import { DictsActions, dictsReducer } from './common/dicts'
import { Section, SectionsActions, sectionsReducer } from './common/sections'
import { ArticleState, ArticleActions, articleReducer } from './pages/article/article'
import { DictArticlesState, DictArticlesActions, dictArticlesReducer } from './pages/dict/dict'
import { LetterFilterState, LetterFilterActions, letterFilterReducer } from './pages/dict/letterfilter'
import { AbbrActions, AbbrState, abbrReducer } from './pages/dict/abbr'
import { PrefaceActions, PrefaceState, prefaceReducer } from './pages/dict/preface'
import { loadingBarReducer } from 'react-redux-loading-bar'

export const rootReducer = combineReducers({
    dicts: dictsReducer,
    sections: sectionsReducer,
    search: searchReducer,
    article: articleReducer,
    letterFilter: letterFilterReducer,
    dictArticles: dictArticlesReducer,
    abbr: abbrReducer,
    preface: prefaceReducer,
    loadingBar: loadingBarReducer,
})

export type RootState = ReturnType<typeof rootReducer>

export function useSelector<TSelected = unknown>(
    selector: (state: RootState) => TSelected,
    equalityFn?: (left: TSelected, right: TSelected) => boolean
): TSelected {
    return useSelectorParent<RootState,TSelected>(selector, equalityFn)
}

export function useDicts(): Dict[] {
    return useSelector<Dict[]>(state => state.dicts)
}

export function useDictsInSection(sectionID: string): Dict[] {
    const section = useSection(sectionID)
    const dicts = useDicts()
    return section.DictIDs.map(dictID => dicts.find(d => d.ID == dictID))
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

export function useSections(): Section[] {
    return useSelector<Section[]>(state => state.sections)
}

export function useSection(id: string): Section | undefined {
    return useSections().find(s => s.ID === id)
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

export function useAbbr(): AbbrState {
    return useSelector<AbbrState>(state => state.abbr)
}

export function usePreface(): PrefaceState {
    return useSelector<PrefaceState>(state => state.preface)
}

export type { AppThunkAction } from './thunk'
