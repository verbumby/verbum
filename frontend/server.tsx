import * as React from 'react'
import { renderToString } from 'react-dom/server'
import Koa from 'koa'

import { App } from './app/App'

const k = new Koa()
k.use(async ctx => {
    ctx.body = renderToString(<App message="ololo" />)
})

console.log('listening on localhost:8079')
k.listen(8079, 'localhost')
