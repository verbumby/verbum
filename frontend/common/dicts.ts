import { Dict } from './dict'
import { Section } from './sections'

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
