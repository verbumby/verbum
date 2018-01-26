import React from 'react'

import Form from './form'

export default class NewPage extends React.Component {
    render() {
        return (<div>
            <h1 className="title">Create New Dictionary</h1>
            <Form onCancelRedirectTo="/dictionaries" />
        </div>)
    }
}
