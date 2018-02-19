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
        return (<div>
            <div className="level">
                <div className="level-left">
                    <h1 className="level-item title">Dictionaries</h1>
                </div>
                <div className="level-right">
                    <div className="field is-grouped level-item">
                        <p class="control">
                            <Link to={`${url}/new`} className="button is-link">New Dictionary</Link>
                        </p>
                    </div>
                </div>
            </div>
            <hr />
            {this.props.data &&
            <table className="table is-hoverable is-fullwidth">
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
