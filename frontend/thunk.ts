import type { ThunkAction } from '@reduxjs/toolkit'

export type AppThunkAction<ReturnType = void> = ThunkAction<Promise<ReturnType>, any, unknown, any>
