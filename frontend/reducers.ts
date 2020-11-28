import { combineReducers, Dispatch } from 'redux'
import { useSelector as useSelectorParent } from 'react-redux'

// TODO: split this file

export type Dictionary = {
    ID: string
    Title: string
}

export type DictionariesListState = Dictionary[]

const DICTIONARIES_LIST_SET = 'DICTIONARIES_LIST/SET'

type DictionariesListSetAction = {
    type: typeof DICTIONARIES_LIST_SET
    dictionaries: DictionariesListState
}

function dictionariesListSet(dictionaries: Dictionary[]): DictionariesListSetAction {
    return {
        type: DICTIONARIES_LIST_SET,
        dictionaries,
    }
}

type DictionariesListActions = DictionariesListSetAction

function dictionariesListReducer(state:DictionariesListState = [], a:DictionariesListActions): DictionariesListState {
    switch (a.type) {
        case DICTIONARIES_LIST_SET:
            return [...a.dictionaries]
        default:
            return state
    }
}

export type Article = {
    ID: string
    Content: string
    DictionaryID: string
}

type SearchState = {
    q: string
    hits: Article[]
}

const SEARCH_KICKOFF = 'SEARCH/KICKOFF'
const SEARCH_SET_HITS = 'SEARCH/SET_HITS'
const SEARCH_RESET = 'SEARCH/RESET'

type SearchKickOffAction = {
    type: typeof SEARCH_KICKOFF,
    q: string
}

function searchKickOff(q: string): SearchKickOffAction {
    return {
        type: SEARCH_KICKOFF,
        q,
    }
}

type SearchSetHitsAction = {
    type: typeof SEARCH_SET_HITS
    hits: Article[]
}

function searchSetHits(hits: Article[]): SearchSetHitsAction {
    return {
        type: SEARCH_SET_HITS,
        hits,
    }
}

type SearchResetAction = {
    type: typeof SEARCH_RESET,
}
export function searchReset(): SearchResetAction {
    return {type: SEARCH_RESET}
}

type SearchActions = SearchKickOffAction | SearchSetHitsAction | SearchResetAction

function searchReducer(state: SearchState = {q: '', hits: []}, a: SearchActions): SearchState {
    switch (a.type) {
        case SEARCH_KICKOFF:
            return {
                q: a.q,
                hits: [],
            }
        case SEARCH_SET_HITS:
            return {
                ...state,
                hits: a.hits,
            }
        case SEARCH_RESET:
            return {
                q: '',
                hits: [],
            }
        default:
            return state
    }
}

type AllActions = DictionariesListActions | SearchActions

export type RootState = {
    dictionaries: DictionariesListState
    search: SearchState
}

export const rootReducer = combineReducers<RootState, AllActions>({
    dictionaries: dictionariesListReducer,
    search: searchReducer,
})

export function useSelector<TSelected = unknown>(
    selector: (state: RootState) => TSelected,
    equalityFn?: (left: TSelected, right: TSelected) => boolean
): TSelected {
    return useSelectorParent<RootState,TSelected>(selector, equalityFn)
}

export function useDicts(): DictionariesListState {
    return useSelector<DictionariesListState>(state => state.dictionaries)
}

export const dictionariesListFetch = () => {
    return async (dispatch: Dispatch) => {
        try {
            dispatch(dictionariesListSet(await verbumClient.getDictionaries()))
        } catch (err) {
            console.log("ERROR: ", err)
            throw err
        }
    }
}

export function useSearchState(): SearchState {
    return useSelector<SearchState>(state => state.search)
}

export const search = (q: string) => {
    return async (dispatch: Dispatch) => {
        try {
            dispatch(searchKickOff(q))
            dispatch(searchSetHits(await verbumClient.search(q)))
        } catch (err) {
            console.log('ERROR: ', err)
            throw err
        }
    }
}
