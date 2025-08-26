
export type Section = {
    ID: string,
    Name: string,
    DictIDs: string[],
    Descr: string,
}

const SECTIONS_SET = 'SECTIONS/SET'
type SectionsSetAction = {
    type: typeof SECTIONS_SET
    sections: Section[]
}

export function sectionsSet(sections: Section[]): SectionsSetAction {
    return { type: SECTIONS_SET, sections }
}

export type SectionsActions = SectionsSetAction

export function sectionsReducer(state: Section[] = [], a: SectionsActions): Section[] {
    switch (a.type) {
        case SECTIONS_SET:
            return [...a.sections]
        default:
            return state
    }
}
