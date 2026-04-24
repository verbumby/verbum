import { createSlice, type PayloadAction } from '@reduxjs/toolkit'
import type { Dict } from './dict'
import type { Section } from './sections'

const dictsSlice = createSlice({
    name: 'dicts',
    initialState: [] as Dict[],
    reducers: {
        dictsSet: (_, action: PayloadAction<Dict[]>) => [...action.payload],
    },
})

export const { dictsSet } = dictsSlice.actions
export const dictsReducer = dictsSlice.reducer

export type DictsMetadata = {
    Dicts: Dict[]
    Sections: Section[]
}
