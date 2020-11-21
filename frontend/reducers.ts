import { combineReducers, Dispatch } from 'redux'
import { useSelector as useSelectorParent } from 'react-redux'

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

function dictionariesListReducer(state:DictionariesListState = [], action:DictionariesListActions): DictionariesListState {
    switch (action.type) {
        case DICTIONARIES_LIST_SET:
            return [...action.dictionaries]
        default:
            return state
    }
}

type AllActions = DictionariesListActions

export type RootState = {
    dictionaries: DictionariesListState
}

export const rootReducer = combineReducers<RootState, AllActions>({
    dictionaries: dictionariesListReducer,
})

export function useSelector<TSelected = unknown>(
    selector: (state: RootState) => TSelected,
    equalityFn?: (left: TSelected, right: TSelected) => boolean
): TSelected {
    return useSelectorParent<RootState,TSelected>(selector, equalityFn)
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
