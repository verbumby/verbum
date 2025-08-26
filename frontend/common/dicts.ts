import { match } from 'react-router-dom'
import { AppThunkAction } from '../store'
import { Dict } from './dict'
import { Section, sectionsSet } from './sections'

const DICTS_SET = 'DICTS/SET'

type DictsSetAction = {
    type: typeof DICTS_SET
    dicts: Dict[]
}

export function dictsSet(dicts: Dict[]): DictsSetAction {
    return { type: DICTS_SET, dicts }
}

export type DictsActions = DictsSetAction

export function dictsReducer(state: Dict[] = [], a: DictsActions): Dict[] {
    switch (a.type) {
        case DICTS_SET:
            return [...a.dicts]
        default:
            return state
    }
}

export type DictsMetadata = {
    Dicts: Dict[],
    Sections: Section[],
}

export const dictsFetch = (match: match, urlSearch: URLSearchParams): AppThunkAction => {
    return async (dispatch) => {
        try {
            const dm = await verbumClient.getDictionaries()
            dispatch(dictsSet(dm.Dicts))
            dispatch(sectionsSet(dm.Sections))
        } catch (err) {
            console.log("ERROR: ", err)
            throw err
        }
    }
}
