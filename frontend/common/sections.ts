import { createSlice, type PayloadAction } from '@reduxjs/toolkit'

export type Section = {
    ID: string
    Name: string
    DictIDs: string[]
    Descr: string
}

const sectionsSlice = createSlice({
    name: 'sections',
    initialState: [] as Section[],
    reducers: {
        sectionsSet: (_, action: PayloadAction<Section[]>) => [
            ...action.payload,
        ],
    },
})

export const { sectionsSet } = sectionsSlice.actions
export const sectionsReducer = sectionsSlice.reducer
