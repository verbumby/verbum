import * as React from 'react'
import { renderToString } from 'react-dom/server'

import { App } from './app/App'

console.log(renderToString(<App message="ololo" />))
