import 'source-map-support/register'

import * as React from 'react'
import { renderToString } from 'react-dom/server'
import Koa from 'koa'
import koaStatic from 'koa-static'
import koaMount from 'koa-mount'

import { App } from './app/App'

const kstatics = new Koa()
kstatics.use(koaStatic('public'))

const k = new Koa()
k.use(koaMount('/statics', kstatics))
k.use(async ctx => {
    ctx.body = renderToString(<App message="ololo" />)
})

console.log('listening on localhost:8079')
k.listen(8079, 'localhost')
