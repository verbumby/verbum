import { Abbrevs, Article, ArticleList, DictsMetadata, LetterFilter, SearchResult, Suggestion } from '../common'

export interface VerbumAPIClient {
    getPreface(dictID: string): Promise<string>
    getAbbr(dictID: string): any
    withSignal(signal: AbortSignal): this
    getDictionaries(): Promise<DictsMetadata>
    search(q: string, inDicts: string, page: number): Promise<SearchResult>
    suggest(q: string, inDicts: string): Promise<Suggestion[]>
    getArticle(dictID: string, articleID: string): Promise<Article>
    getLetterFilter(dictID: string, prefix: string): Promise<LetterFilter>
    getDictArticles(dictID: string, q: string, prefix: string, page: number): Promise<ArticleList>
    getIndexHTML(): Promise<string>
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
    abstract callString(path: string): Promise<string>

    async getDictionaries(): Promise<DictsMetadata> {
        return this.call<DictsMetadata>('/api/dictionaries')
    }

    async search(q: string, inDicts: string, page: number): Promise<SearchResult> {
        q = encodeURIComponent(q)
        inDicts = encodeURIComponent(inDicts)
        return this.call<SearchResult>(`/api/search?q=${q}&in=${inDicts}&page=${page}`)
    }

    async suggest(q: string, inDicts: string): Promise<Suggestion[]> {
        q = encodeURIComponent(q)
        inDicts = encodeURIComponent(inDicts)
        return this.call<Suggestion[]>(`/api/suggest?q=${q}&in=${inDicts}`)
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

    async getDictArticles(dictID: string, q: string, prefix: string, page: number): Promise<ArticleList> {
        dictID = encodeURIComponent(dictID)
        q = encodeURIComponent(q)
        prefix = encodeURIComponent(prefix)
        return this.call<ArticleList>(`/api/search?q=${q}&in=${dictID}&prefix=${prefix}&page=${page}&track_total_hits=true`)
    }

    async getAbbr(dictID: string): Promise<Abbrevs> {
        dictID = encodeURIComponent(dictID)
        return this.call<Abbrevs>(`/api/dictionaries/${dictID}/abbrevs`)
    }

    async getPreface(dictID: string): Promise<string> {
        dictID = encodeURIComponent(dictID)
        return this.call<string>(`/api/dictionaries/${dictID}/preface`)
    }

    async getIndexHTML(): Promise<string> {
        return this.callString(`/api/index.html`)
    }
}
