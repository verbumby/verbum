import { AppThunkAction } from '../store'
import { Dict } from './dict'

export type DictsState = Dict[]

const DICTS_SET = 'DICTS/SET'
type DictsSetAction = {
    type: typeof DICTS_SET
    dicts: DictsState
}
function dictsSet(dicts: Dict[]): DictsSetAction {
    return { type: DICTS_SET, dicts }
}

export type DictsActions = DictsSetAction

export function dictsReducer(state:DictsState = [], a:DictsActions): DictsState {
    switch (a.type) {
        case DICTS_SET:
            return [...a.dicts]
        default:
            return state
    }
}

export const dictsFetch = (urlSearch: URLSearchParams): AppThunkAction => {
    return async (dispatch) => {
        try {
            dispatch(dictsSet(await verbumClient.getDictionaries()))
        } catch (err) {
            console.log("ERROR: ", err)
            throw err
        }
    }
}
