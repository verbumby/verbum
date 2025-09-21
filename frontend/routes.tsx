import * as React from 'react'
import { articleFetchServer, ArticlePage } from './pages/article'
import { IndexPage, searchServer } from './pages/index/index'
import { AppThunkAction } from './store'
import { DictPage, letterFilterFetchServer, dictArticlesFetchServer, abbrFetchServer, prefaceFetchServer } from './pages/dict'
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
