import { match } from "react-router-dom"
import { Article } from "../../common"
import { AppThunkAction } from "../../store"

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
function articleFetchSuccess(a: Article): ArticleFetchSuccessAction {
    return {
        type: ARTICLE_FETCH_SUCCESS,
        a,
    }
}

const ARTICLE_RESET = 'ARTICLE/RESET'
type ArticleResetAction = {
    type: typeof ARTICLE_RESET,
}
export function articleReset(): ArticleResetAction {
    return { type: ARTICLE_RESET }
}

export type ArticleActions = ArticleFetchKickOffAction | ArticleFetchSuccessAction | ArticleResetAction

export function articleReducer(state: Article = null, a: ArticleActions): Article {
    switch (a.type) {
        case ARTICLE_FETCH_KICKOFF:
            return null
        case ARTICLE_FETCH_SUCCESS:
            return a.a
        case ARTICLE_RESET:
            return null
        default:
            return state
    }
}

export const articleFetch = (match: match<MatchParams>, urlSearch: URLSearchParams): AppThunkAction => {
    return async (dispatch, getState): Promise<void> => {
        try {
            const { dictID, articleID } = match.params
            dispatch(articleFetchKickOff(dictID, articleID))
            dispatch(articleFetchSuccess(await verbumClient.getArticle(dictID, articleID)))
        } catch (err) {
            console.log('ERROR: ', err)
            throw err
        }
    }
}
