import 'source-map-support/register'

import { readFileSync } from 'fs'
import * as React from 'react'
import { renderToString } from 'react-dom/server'
import Koa from 'koa'
import koaStatic from 'koa-static'
import koaMount from 'koa-mount'
import { matchPath, StaticRouter, StaticRouterContext } from 'react-router'
import { configureStore } from '@reduxjs/toolkit'
import { Provider } from 'react-redux'
import { Helmet } from "react-helmet";

import { App } from './App'
import { rootReducer } from './store'
import { VerbumAPIClientServer } from './verbum/server'
import { routes } from './routes'

global.verbumClient = new VerbumAPIClientServer({apiURL: 'http://localhost:8080'})

const indexhtml = readFileSync('index.html', 'utf-8')

const kstatics = new Koa()
kstatics.use(koaStatic(
    'public',
    {
        maxage: 1E3 * 60 * 60 * 24 * 30,
        immutable: true,
    },
))

const k = new Koa()
k.use(koaMount('/statics', kstatics))
k.use(async ctx => {
    const store = configureStore({
        reducer: rootReducer,
    })
    const promises: Promise<void>[] = []
    routes.some(route => {
        const match = matchPath(ctx.URL.pathname, route)
        if (match) {
            for (const dataLoader of route.dataLoaders) {
                promises.push(store.dispatch(dataLoader(match, ctx.URL.searchParams)))
            }
        }
        return match
    })
    await Promise.all(promises)

    const preloadedState = store.getState()

    const routerContext: StaticRouterContext = {}
    const reactRendered = renderToString(
        <Provider store={store}>
            <StaticRouter location={ctx.url} context={routerContext}>
                <App />
            </StaticRouter>
        </Provider>
    )
    const helmet = Helmet.renderStatic()

    if (routerContext.url) {
        ctx.redirect(routerContext.url)
        return
    }

    let body = indexhtml
    body = body.replace('HEAD_TITLE_PLACEHOLDER', helmet.title.toString())
    body = body.replace('HEAD_META_PLACEHOLDER', helmet.meta.toString())
    body = body.replace('PRELOADED_STATE_PLACEHOLDER', JSON.stringify(preloadedState).replace(/</g, '\\u003c'))
    body = body.replace('BODY_PLACEHOLDER', reactRendered)
    ctx.body = body
})

console.log('listening on localhost:8079')
k.listen(8079, 'localhost')
