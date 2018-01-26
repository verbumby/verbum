import React from 'react'
import { connect } from 'react-redux'
import { bindActionCreators } from 'redux'

import Form from './form'
import { createDictionary } from '../actions'
import { ifOK } from '../../utils'

class NewPage extends React.Component {
    render() {
        const onSave = (data) => {
            this.props.createDictionary(data)
                .then(ifOK((data) => this.props.history.push('/dictionaries')))
        }
        return (<div>
            <h1 className="title">Create New Dictionary</h1>
            <Form
                formData={{}}
                onSave={onSave}
                onCancelRedirectTo="/dictionaries"
            />
        </div>)
    }
}

export default connect(
    state => {},
    dispatch => bindActionCreators({ createDictionary }, dispatch),
)(NewPage)
