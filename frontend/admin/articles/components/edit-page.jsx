import React from 'react'
import { connect } from 'react-redux'
import { bindActionCreators } from 'redux'

import { fetchRecord, leaveRecord} from '../actions'
import Form from './form'

class EditPage extends React.Component {
    componentWillMount() {
        this.props.fetchRecord(this.props.match.params.articleID)
    }

    render() {
        if (!this.props.data) {
            return null
        }

        console.log(this.props.data)
        return (<div>
            <h1 className="title">Edit Article #{this.props.data.id}</h1>
            <Form
                formData={this.props.data}
                onCancelRedirectTo="/articles"
            />
        </div>)
    }
}

export default connect(
    state => ({
        data: state.articles.record.data,
    }),
    dispatch => bindActionCreators({ fetchRecord, leaveRecord }, dispatch)
)(EditPage)
