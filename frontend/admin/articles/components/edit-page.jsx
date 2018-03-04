import React from 'react'
import { connect } from 'react-redux'
import { bindActionCreators } from 'redux'

import { fetchRecord, leaveRecord, updateRecord} from '../actions'
import Form from './form'

class EditPage extends React.Component {
    componentWillMount() {
        this.props.fetchRecord(this.props.match.params.articleID)
    }

    componentWillUnmount() {
        this.props.leaveRecord()
    }

    render() {
        if (!this.props.data) {
            return null
        }

        const onSave = (data) => {
            this.props.updateRecord(data)
        }

        return (<div>
            <h1 className="title">Edit #{this.props.data.ID} `{this.props.data.Title}`</h1>
            <Form
                formData={this.props.data}
                onSave={onSave}
                onCancelRedirectTo="/articles"
            />
        </div>)
    }
}

export default connect(
    state => ({
        data: state.articles.record.data,
    }),
    dispatch => bindActionCreators({ fetchRecord, leaveRecord, updateRecord }, dispatch)
)(EditPage)
