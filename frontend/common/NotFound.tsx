import * as React from 'react'
import { Route } from 'react-router-dom'

export const NotFound: React.VFC = () => (
    <Route render={({ staticContext }) => {
        if (staticContext) staticContext.statusCode = 404;
        return (<div>Такой старонкі не існуе.</div>)
    }} />
)
