import type { Article } from './article'
import type { Pagination } from './pagination'

export type SearchResult = {
    Articles: Article[]
    TermSuggestions: string[]
    Pagination: Pagination
}
