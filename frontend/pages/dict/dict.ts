import { match } from "react-router-dom"
import { ArticleList, URLSearch, useURLSearch as useURLSearchCommon } from "../../common"
import { AppThunkAction } from "../../store"

export type MatchParams = {
    dictID: string
}

export const URLSearchDefaults = {
    prefix: '',
    page: 1,
}

export const useURLSearch = () => useURLSearchCommon(URLSearchDefaults)

export type DictArticlesState = ArticleList

const DICT_ARTICLES_FETCH_KICKOFF = 'DICT_ARTICLES/FETCH/KICKOFF'
type DictArticlesFetchKickoffAction = {
    type: typeof DICT_ARTICLES_FETCH_KICKOFF
    dictID: string
    prefix: string
    page: number
}
function dictArticlesFetchKickOff(dictID: string, prefix: string, page: number): DictArticlesFetchKickoffAction {
    return { type: DICT_ARTICLES_FETCH_KICKOFF, dictID, prefix, page }
}

const DICT_ARTICLES_FETCH_SUCCESS = 'DICT_ARTICLES/FETCH/SUCCESS'
type DictArticlesFetchSuccessAction = {
    type: typeof DICT_ARTICLES_FETCH_SUCCESS
    articleList: ArticleList
}
function dictArticlesFetchSuccess(articleList: ArticleList): DictArticlesFetchSuccessAction {
    return { type: DICT_ARTICLES_FETCH_SUCCESS, articleList }
}

const DICT_ARTICLES_FETCH_FAILURE = 'DICT_ARTICLES/FETCH/FAILURE'
type DictArticlesFetchFailureAction = {
    type: typeof DICT_ARTICLES_FETCH_FAILURE
}
function dictArticlesFetchFailure(): DictArticlesFetchFailureAction {
    return { type: DICT_ARTICLES_FETCH_FAILURE }
}

const DICT_ARTICLES_RESET = 'DICT_ARTICLES/RESET'
type DictArticlesResetAction = {
    type: typeof DICT_ARTICLES_RESET
}
export function dictArticlesReset(): DictArticlesResetAction {
    return { type: DICT_ARTICLES_RESET }
}

export type DictArticlesActions = DictArticlesFetchKickoffAction | DictArticlesFetchSuccessAction | DictArticlesFetchFailureAction | DictArticlesResetAction

export function dictArticlesReducer(state: DictArticlesState = null, a: DictArticlesActions): DictArticlesState {
    switch (a.type) {
        case DICT_ARTICLES_FETCH_KICKOFF:
            return null
        case DICT_ARTICLES_FETCH_SUCCESS:
            return a.articleList
        case DICT_ARTICLES_RESET:
            return null
        default:
            return state
    }
}

export const dictArticlesFetch = (match: match<MatchParams>, urlSearch: URLSearch<typeof URLSearchDefaults>): AppThunkAction => {
    return async (dispatch, getState): Promise<void> => {
        try {
            const { dictID } = match.params
            const prefix = urlSearch.get('prefix')
            const page = urlSearch.get('page')
            const state = getState()
            if (state.dictArticles
                && state.dictArticles.DictID === dictID
                && state.dictArticles.Prefix === prefix
                && state.dictArticles.Pagination.Current === page
            ) {
                return
            }
            dispatch(dictArticlesFetchKickOff(dictID, prefix, page))
            dispatch(dictArticlesFetchSuccess(await verbumClient.getDictArticles(dictID, prefix, page)))
        } catch (err) {
            dispatch(dictArticlesFetchFailure())
            console.log('ERROR: ', err)
            throw err
        }
    }
}

export const dictArticlesFetchServer = (match: match<MatchParams>, urlSearchParams: URLSearchParams): AppThunkAction =>
    dictArticlesFetch(match, new URLSearch(URLSearchDefaults, urlSearchParams))
