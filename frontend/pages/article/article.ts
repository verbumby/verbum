import { createSlice, type PayloadAction } from '@reduxjs/toolkit'
import type { Article } from '../../common/article'
import type { AppThunkAction } from '../../thunk'

export type ArticleState = {
    a?: Article
}

export type MatchParams = {
    dictID: string
    articleID: string
}

const articleSlice = createSlice({
    name: 'article',
    initialState: {} as ArticleState,
    reducers: {
        articleFetchKickOff: (state) => state,
        articleFetchSuccess: (_, action: PayloadAction<Article>) => ({
            a: action.payload,
        }),
        articleFetchFailure: (state) => state,
        articleReset: () => ({}),
    },
})

const { articleFetchKickOff, articleFetchFailure } = articleSlice.actions
export const { articleFetchSuccess, articleReset } = articleSlice.actions
export const articleReducer = articleSlice.reducer

export const articleFetch = (params: Partial<MatchParams>): AppThunkAction => {
    return async (dispatch, getState): Promise<void> => {
        try {
            const { dictID, articleID } = params
            dispatch(articleFetchKickOff(dictID, articleID))
            dispatch(
                articleFetchSuccess(
                    await verbumClient.getArticle(dictID, articleID),
                ),
            )
        } catch (err) {
            dispatch(articleFetchFailure())
            console.log('ERROR: ', err)
            throw err
        }
    }
}

export const articleFetchServer = (
    params: MatchParams,
    urlSearchParams: URLSearchParams,
): AppThunkAction => articleFetch(params)
