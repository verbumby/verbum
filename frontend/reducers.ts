import { combineReducers } from 'redux'
import { useSelector as useSelectorParent } from 'react-redux'

type Dictionary = {
    ID: string
    Name: string
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

export const rootReducer = combineReducers({
    dictionaries: dictionariesListReducer,
})

export type RootState = ReturnType<typeof rootReducer>

export function useSelector<TSelected = unknown>(
    selector: (state: RootState) => TSelected,
    equalityFn?: (left: TSelected, right: TSelected) => boolean
): TSelected {
    return useSelectorParent<RootState,TSelected>(selector, equalityFn)
}
