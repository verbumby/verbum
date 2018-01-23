import React from 'react'
import ReactDOM from 'react-dom'
import { Provider } from 'react-redux'
import { createStore, combineReducers, applyMiddleware } from 'redux'
import thunkMiddleware from 'redux-thunk'
import { BrowserRouter as Router, Switch, Route } from 'react-router-dom'

import TopNav from './core/components/top-nav'
import HomePage from './home/components/home-page'
import Messages from './messages/messages'

import messagesReducer from './messages/reducers'

const rootReducer = combineReducers({
    user: (state = {}) => state,
    config: (state = {}) => state,
    messages: messagesReducer,
})

const store = createStore(
    rootReducer,
    window.verbumInitData,
    applyMiddleware(
        thunkMiddleware,
    ),
)

ReactDOM.render(
    <Provider store={store}>
        <Router basename="/admin">
            <div>
                <div>
                    <TopNav />
                    <Messages />
                </div>
                <div className="container-fluid">
                    <Switch>
                        <Route exact path="/" component={HomePage} />
                    </Switch>
                </div>
            </div>
        </Router>
    </Provider>,
    document.getElementById('root')
)
