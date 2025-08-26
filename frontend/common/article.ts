import { Pagination } from "./"

export type Article = {
    ID: string
    Title: string
    Headword: string[]
    Content: string
    DictionaryID: string
}

export type ArticleList = {
    DictIDs: string[]
    Q: string
    Prefix: string
    Articles: Article[]
    TermSuggestions: string[]
    Pagination: Pagination
}
