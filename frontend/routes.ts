import * as React from 'react'
import { match } from 'react-router-dom'
import { articleFetch, ArticlePage } from './pages/article'
import { IndexPage, search } from './pages/index/index'
import { AppThunkAction } from './store'
import { dictsFetch } from './common'

type DataLoader = (match: match, urlSearch: URLSearchParams) => AppThunkAction

type Route = {
    path: string,
    component: React.ComponentType<any>,
    dataLoaders: DataLoader[],
}

export const routes: Route[] = [
    {
        path: '/:dictID/:articleID',
        component: ArticlePage,
        dataLoaders: [dictsFetch, articleFetch],
    },
    {
        path: '/',
        component: IndexPage,
        dataLoaders: [dictsFetch, search],
    },
]