import * as React from 'react'
import { ArticlePage } from './pages/article/ArticlePage'
import { articleFetchServer } from './pages/article/article'
import { IndexPage } from './pages/index/IndexPage'
import { searchServer } from './pages/index/search'
import { AppThunkAction } from './store'
import { DictPage } from './pages/dict/DictPage'
import { letterFilterFetchServer } from './pages/dict/letterfilter'
import { dictArticlesFetchServer } from './pages/dict/dict'
import { abbrFetchServer } from './pages/dict/abbr'
import { prefaceFetchServer } from './pages/dict/preface'
import { SupportPage } from './pages/support/SupportPage'

type DataLoader = (match: {}, urlSearch: URLSearchParams) => AppThunkAction

type Route = {
    path: string,
    element: React.ReactElement,
    dataLoaders: DataLoader[],
}

export const routes: Route[] = [
    {
        path: '/s/:sectionID',
        element: <IndexPage />,
        dataLoaders: [searchServer],
    },
    {
        path: '/:dictID/:articleID',
        element: <ArticlePage />,
        dataLoaders: [articleFetchServer],
    },
    {
        path: '/support',
        element: <SupportPage />,
        dataLoaders: [],
    },
    {
        path: '/:dictID',
        element: <DictPage />,
        dataLoaders: [letterFilterFetchServer, dictArticlesFetchServer, abbrFetchServer, prefaceFetchServer],
    },
    {
        path: '/',
        element: <IndexPage />,
        dataLoaders: [searchServer],
    },
]
