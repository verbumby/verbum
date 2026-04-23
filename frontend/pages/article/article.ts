import { createSlice, type PayloadAction } from '@reduxjs/toolkit'
import type { Article } from '../../common/article'
import { serverLoader } from '../../common/serverLoader'
import { URLSearch } from '../../common/urlsearch'
import type { AppThunkAction } from '../../thunk'

export type ArticleState = {
    a?: Article | null
}

export type MatchParams = {
    dictID: string
    articleID: string
}

export const URLSearchDefaults = {}

const articleSlice = createSlice({
    name: 'article',
    initialState: {} as ArticleState,
    reducers: {
        articleFetchKickOff: (state) => state,
        articleFetchSuccess: (_, action: PayloadAction<Article | null>) => ({
            a: action.payload,
        }),
        articleFetchFailure: (state) => state,
        articleReset: () => ({}),
    },
})

const { articleFetchKickOff, articleFetchFailure } = articleSlice.actions
export const { articleFetchSuccess, articleReset } = articleSlice.actions
export const articleReducer = articleSlice.reducer

export const articleFetch = (
    params: MatchParams,
    _urlSearch?: URLSearch<typeof URLSearchDefaults>,
): AppThunkAction => {
    return async (dispatch, _getState): Promise<void> => {
        try {
            const { dictID, articleID } = params
            dispatch(articleFetchKickOff())
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

export const articleFetchServer = serverLoader(URLSearchDefaults, articleFetch)
