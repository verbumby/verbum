import React from 'react'
import { connect } from 'react-redux'
import { bindActionCreators } from 'redux'

import { fetchList } from '../../../tasks/actions'
import { ifOK } from '../../../utils'

class FilterTask extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            tasks: null
        }
    }
    componentWillMount() {
        this.props.fetchList({}).then(ifOK( data => this.setState({ tasks: data.Data }) ))
    }

    render() {
        if (!this.state.tasks) {
            return null
        }

        const { value, onChange } = this.props
        const tasks = this.state.tasks
        return (<div class="select">
            <select value={value} onChange={(ev) => onChange(parseInt(ev.target.value))}>
                <option value="-1">- Filter by Task -</option>
                {tasks.map(d => <option value={d.ID}>{d.Title}</option>)}
            </select>
        </div>)
    }
}

export default connect(
    state => ({}),
    dispatch => bindActionCreators({ fetchList }, dispatch),

)(FilterTask)
