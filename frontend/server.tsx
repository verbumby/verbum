import 'source-map-support/register'

import { readFileSync } from 'fs'
import * as React from 'react'
import { renderToString } from 'react-dom/server'
import Koa from 'koa'
import koaStatic from 'koa-static'
import koaMount from 'koa-mount'
import { StaticRouter, StaticRouterContext } from 'react-router'
import { configureStore } from '@reduxjs/toolkit'
import { Provider } from 'react-redux'

import { App } from './app/App'
import { dictionariesListFetch, rootReducer } from './reducers'
import { VerbumAPIClientServer } from './verbum/server'

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
    const stateStore = configureStore({
        reducer: rootReducer,
    })
    await stateStore.dispatch(dictionariesListFetch())
    const preloadedState = stateStore.getState()

    const routerContext: StaticRouterContext = {}
    const reactRendered = renderToString(
        <Provider store={stateStore}>
            <StaticRouter location={ctx.url} context={routerContext}>
                <App />
            </StaticRouter>
        </Provider>
    )

    if (routerContext.url) {
        ctx.redirect(routerContext.url)
        return
    }


    let body = indexhtml
    body = body.replace('HEAD_TITLE_PLACEHOLDER', 'Some TItle ololo')
    body = body.replace('PRELOADED_STATE_PLACEHOLDER', JSON.stringify(preloadedState).replace(/</g, '\\u003c'))
    body = body.replace('BODY_PLACEHOLDER', reactRendered)
    ctx.body = body
})

console.log('listening on localhost:8079')
k.listen(8079, 'localhost')
