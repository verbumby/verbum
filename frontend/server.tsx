import 'source-map-support/register'

import { readFileSync } from 'fs'
import * as React from 'react'
import { renderToString } from 'react-dom/server'
import Koa from 'koa'
import koaStatic from 'koa-static'
import koaMount from 'koa-mount'
import { StaticRouter, StaticRouterContext } from 'react-router'
import { createStore, Reducer } from 'redux'
import { Provider } from 'react-redux'
import { devToolsEnhancer } from 'redux-devtools-extension'

import { App } from './app/App'
import { rootReducer, RootState } from './reducers'

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
    const stateStore = createStore(rootReducer, devToolsEnhancer({}))
    const routerContext: StaticRouterContext = {}
    const reactRendered = renderToString(
        <Provider store={stateStore}>
            <StaticRouter location={ctx.url} context={routerContext}>
                <App message="ololo" />
            </StaticRouter>
        </Provider>
    )

    if (routerContext.url) {
        ctx.redirect(routerContext.url)
        return
    }

    const preloadedState = stateStore.getState()

    let body = indexhtml
    body = body.replace('HEAD_TITLE_PLACEHOLDER', 'Some TItle ololo')
    body = body.replace('PRELOADED_STATE_PLACEHOLDER', JSON.stringify(preloadedState).replace(/</g, '\\u003c'))
    body = body.replace('BODY_PLACEHOLDER', reactRendered)
    ctx.body = body
})

console.log('listening on localhost:8079')
k.listen(8079, 'localhost')
