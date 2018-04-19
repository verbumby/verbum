import React from 'react'
import { Link } from 'react-router-dom'
import { connect } from 'react-redux'
import { bindActionCreators } from 'redux'

import { fetchList, leaveList } from '../actions'

class ListPage extends React.Component {
    componentWillMount() {
        this.props.fetchList()
    }

    componentWillUnmount() {
        this.props.leaveList()
    }

    render() {
        const { url } = this.props.match
        return (<div className="mt-2">
            <h1 className="d-inline-block mr-2 align-middle">Dictionaries</h1>
            <Link to={`${url}/new`} className="btn btn-light align-middle">New Dictionary</Link>

            {this.props.data &&
                <table className="table">
                    <thead>
                        <tr><th>Title</th></tr>
                    </thead>
                    <tbody>
                        {this.props.data.map(item => <tr><td>{item.Title}</td></tr>)}
                    </tbody>
                </table>
            }

        </div>)
    }
}

export default connect(
    state => ({
        data: state.dictionaries.list,
    }),
    dispatch => bindActionCreators({ fetchList, leaveList }, dispatch),
)(ListPage)
