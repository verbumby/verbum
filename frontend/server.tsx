import 'source-map-support/register'

import { readFileSync } from 'fs'
import * as React from 'react'
import { renderToString } from 'react-dom/server'
import Koa from 'koa'
import koaStatic from 'koa-static'
import koaMount from 'koa-mount'
import { StaticRouter, StaticRouterContext } from 'react-router'

import { App } from './app/App'

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
    const routerContext: StaticRouterContext = {}
    const reactRendered = renderToString(
        <StaticRouter location={ctx.url} context={routerContext}>
            <App message="ololo" />
        </StaticRouter>
    )

    if (routerContext.url) {
        ctx.redirect(routerContext.url)
        return
    }

    let body = indexhtml
    body = body.replace('HEAD_TITLE_PLACEHOLDER', 'Some TItle ololo')
    body = body.replace('BODY_PLACEHOLDER', reactRendered)
    ctx.body = body
})

console.log('listening on localhost:8079')
k.listen(8079, 'localhost')
