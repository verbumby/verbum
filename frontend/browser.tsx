import * as React from 'react'
import { hydrate } from 'react-dom'

import { App } from './app/App'

hydrate(<App message="ololo" />, document.querySelector('body .content'))
