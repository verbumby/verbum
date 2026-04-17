import { Article } from './article'
import { Pagination } from './pagination'

export type SearchResult = {
    Articles: Article[]
    TermSuggestions: string[]
    Pagination: Pagination
}
