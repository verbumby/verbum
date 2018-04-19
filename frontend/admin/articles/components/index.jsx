import React from 'react'
import { Switch, Route } from 'react-router-dom'

import ListPage from './list-page'
import NewPage from './new-page'
import EditPage from './edit-page'

export default class Index extends React.Component {
    render() {
        const { path } = this.props.match
        return (
            <Switch>
                <Route exact path={path} component={ListPage} />
                <Route path={`${path}/new`} component={NewPage} />
                <Route path={`${path}/:articleID/edit`} component={EditPage} />
            </Switch>
        )
    }
}
