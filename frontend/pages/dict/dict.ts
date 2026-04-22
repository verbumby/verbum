import { createSlice, type PayloadAction } from '@reduxjs/toolkit'
import type { ArticleList } from '../../common/article'
import { useURLSearch as useURLSearchCommon } from '../../common/hooks'
import { serverLoader } from '../../common/serverLoader'
import { URLSearch } from '../../common/urlsearch'
import type { AppThunkAction } from '../../thunk'

export type MatchParams = {
    dictID: string
}

export const URLSearchDefaults = {
    q: '',
    prefix: '',
    page: 1,
    section: '',
}

export const useURLSearch = () => useURLSearchCommon(URLSearchDefaults)

export type DictArticlesState = ArticleList

const dictArticlesSlice = createSlice({
    name: 'dictArticles',
    initialState: null as DictArticlesState | null,
    reducers: {
        dictArticlesFetchKickOff: () => null,
        dictArticlesFetchSuccess: (_, action: PayloadAction<ArticleList>) =>
            action.payload,
        dictArticlesFetchFailure: (state) => state,
        dictArticlesReset: () => null,
    },
})

const { dictArticlesFetchKickOff, dictArticlesFetchFailure } =
    dictArticlesSlice.actions
export const { dictArticlesFetchSuccess, dictArticlesReset } =
    dictArticlesSlice.actions
export const dictArticlesReducer = dictArticlesSlice.reducer

export const dictArticlesFetch = (
    params: Partial<MatchParams>,
    urlSearch: URLSearch<typeof URLSearchDefaults>,
): AppThunkAction => {
    return async (dispatch, getState): Promise<void> => {
        try {
            if (urlSearch.get('section') !== '') {
                return
            }

            const { dictID } = params
            const q = urlSearch.get('q')
            const prefix = urlSearch.get('prefix')
            const page = urlSearch.get('page')
            const state = getState()
            if (
                state.dictArticles &&
                state.dictArticles.DictIDs.length == 1 &&
                state.dictArticles.DictIDs[0] === dictID &&
                state.dictArticles.Q === q &&
                state.dictArticles.Prefix === prefix &&
                state.dictArticles.Pagination.Current === page
            ) {
                return
            }
            dispatch(dictArticlesFetchKickOff())
            dispatch(
                dictArticlesFetchSuccess(
                    await verbumClient.getDictArticles(dictID, q, prefix, page),
                ),
            )
        } catch (err) {
            dispatch(dictArticlesFetchFailure())
            console.log('ERROR: ', err)
            throw err
        }
    }
}

export const dictArticlesFetchServer = serverLoader(URLSearchDefaults, dictArticlesFetch)
