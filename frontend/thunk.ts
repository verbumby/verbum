import type { ThunkAction } from '@reduxjs/toolkit'
import type { RootState } from './store'

export type AppThunkAction<ReturnType = void> = ThunkAction<
    Promise<ReturnType>,
    RootState,
    unknown,
    any
>
