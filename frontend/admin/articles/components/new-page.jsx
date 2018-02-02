import React from 'react'
import { connect } from 'react-redux'
import { bindActionCreators } from 'redux'

import Form from './form'
import { createRecord } from '../actions'
import { ifOK } from '../../utils'

class NewPage extends React.Component {
    render() {
        const onSave = (data) => {
            this.props.createRecord(data)
                .then(ifOK((data) => this.props.history.push('/articles')))
        }
        return (<div>
            <h1 className="title">Create New Article</h1>
            <Form
                formData={{}}
                onSave={onSave}
                onCancelRedirectTo="/articles"
            />
        </div>)
    }
}

export default connect(
    state => {},
    dispatch => bindActionCreators({ createRecord }, dispatch),
)(NewPage)
