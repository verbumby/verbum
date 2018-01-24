import React from 'react'
import { Switch, Route } from 'react-router-dom'

import ListPage from './list-page'
import NewPage from './new-page'

export default class Index extends React.Component {
    render() {
        const { path } = this.props.match
        return (
            <div className="container is-fluid">
                <Switch>
                    <Route exact path={path} component={ListPage} />
                    <Route path={`${path}/new`} component={NewPage} />
                </Switch>
            </div>
        )
    }
}
