import { match } from "react-router-dom"
import { Article } from "../../common"
import { AppThunkAction } from "../../store"

export type ArticleState = {
    a?: Article
}

export type MatchParams = {
    dictID: string
    articleID: string
}

const ARTICLE_FETCH_KICKOFF = 'ARTICLE/FETCH/KICKOFF'
type ArticleFetchKickOffAction = {
    type: typeof ARTICLE_FETCH_KICKOFF,
    dictID: string
    articleID: string
}
function articleFetchKickOff(dictID: string, articleID: string): ArticleFetchKickOffAction {
    return {
        type: ARTICLE_FETCH_KICKOFF,
        dictID,
        articleID,
    }
}

const ARTICLE_FETCH_SUCCESS = 'ARTICLE/FETCH/SUCCESS'
type ArticleFetchSuccessAction = {
    type: typeof ARTICLE_FETCH_SUCCESS,
    a: Article
}
export function articleFetchSuccess(a: Article): ArticleFetchSuccessAction {
    return { type: ARTICLE_FETCH_SUCCESS, a }
}

const ARTICLE_FETCH_FAILURE = 'ARTICLE/FETCH/FAILURE'
type ArticleFetchFailureAction = {
    type: typeof ARTICLE_FETCH_FAILURE
}
function articleFetchFailure(): ArticleFetchFailureAction {
    return { type: ARTICLE_FETCH_FAILURE }
}

const ARTICLE_RESET = 'ARTICLE/RESET'
type ArticleResetAction = {
    type: typeof ARTICLE_RESET,
}
export function articleReset(): ArticleResetAction {
    return { type: ARTICLE_RESET }
}

export type ArticleActions = ArticleFetchKickOffAction | ArticleFetchSuccessAction | ArticleFetchFailureAction | ArticleResetAction

export function articleReducer(state: ArticleState = {}, a: ArticleActions): ArticleState {
    switch (a.type) {
        case ARTICLE_FETCH_KICKOFF:
            return state
        case ARTICLE_FETCH_SUCCESS:
            return { a: a.a }
        case ARTICLE_RESET:
            return {}
        default:
            return state
    }
}

export const articleFetch = (match: match<MatchParams>): AppThunkAction => {
    return async (dispatch, getState): Promise<void> => {
        try {
            const { dictID, articleID } = match.params
            dispatch(articleFetchKickOff(dictID, articleID))
            dispatch(articleFetchSuccess(await verbumClient.getArticle(dictID, articleID)))
        } catch (err) {
            dispatch(articleFetchFailure())
            console.log('ERROR: ', err)
            throw err
        }
    }
}

export const articleFetchServer = (match: match<MatchParams>, urlSearchParams: URLSearchParams): AppThunkAction =>
    articleFetch(match)
