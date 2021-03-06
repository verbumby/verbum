import { Article, ArticleList, Dict, LetterFilter, SearchResult, Suggestion } from '../common'

export interface VerbumAPIClient {
    withSignal(signal: AbortSignal): this
    getDictionaries(): Promise<Dict[]>
    search(q: string, page: number): Promise<SearchResult>
    suggest(q: string): Promise<Suggestion[]>
    getArticle(dictID: string, articleID: string): Promise<Article>
    getLetterFilter(dictID: string, prefix: string): Promise<LetterFilter>
    getDictArticles(dictID: string, prefix: string, page: number): Promise<ArticleList>
}

declare global {
    var verbumClient: VerbumAPIClient
}

export abstract class VerbumAPIClientImpl implements VerbumAPIClient {
    protected signal: AbortSignal;

    withSignal(signal: AbortSignal): this {
        const result = Object.create(this) as this
        result.signal = signal
        return result
    }

    abstract call<T>(path: string): Promise<T>

    async getDictionaries(): Promise<Dict[]> {
        return this.call<Dict[]>('/api/dictionaries')
    }

    async search(q: string, page: number): Promise<SearchResult> {
        q = encodeURIComponent(q)
        return this.call<SearchResult>(`/api/search?q=${q}&page=${page}`)
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
