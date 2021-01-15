import { Article, ArticleList, Dict, LetterFilter, Suggestion } from '../common'

export interface VerbumAPIClient {
    getDictionaries(): Promise<Dict[]>
    search(q: string): Promise<Article[]>
    suggest(q: string): Promise<Suggestion[]>
    getArticle(dictID: string, articleID: string): Promise<Article>
    getLetterFilter(dictID: string, prefix: string): Promise<LetterFilter>
    getDictArticles(dictID: string, prefix: string, page: number): Promise<ArticleList>
}

declare global {
    var verbumClient: VerbumAPIClient
}

export abstract class VerbumAPIClientImpl implements VerbumAPIClient {
    abstract call<T>(path: string): Promise<T>

    async getDictionaries(): Promise<Dict[]> {
        return this.call<Dict[]>('/api/dictionaries')
    }

    async search(q: string): Promise<Article[]> {
        return this.call<Article[]>('/api/search?q=' + encodeURIComponent(q))
    }

    async suggest(q: string): Promise<Suggestion[]> {
        return this.call<Suggestion[]>('/api/suggest?q=' + encodeURIComponent(q))
    }

    async getArticle(dictID: string, articleID: string): Promise<Article> {
        dictID = encodeURIComponent(dictID)
        articleID = encodeURIComponent(articleID)
        return this.call<Article>(`/api/dictionaries/${dictID}/articles/${articleID}`)
    }

    async getLetterFilter(dictID: string, prefix: string): Promise<LetterFilter> {
        dictID = encodeURIComponent(dictID)
        prefix = encodeURIComponent(prefix)
        return this.call<LetterFilter>(`/api/dictionaries/${dictID}/letterfilter?prefix=${prefix}`)
    }

    async getDictArticles(dictID: string, prefix: string, page: number): Promise<ArticleList> {
        dictID = encodeURIComponent(dictID)
        prefix = encodeURIComponent(prefix)
        return this.call<ArticleList>(`/api/dictionaries/${dictID}/articles?prefix=${prefix}&page=${page}`)
    }
}
