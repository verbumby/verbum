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
                    <h1 className="level-item title">Articles</h1>
                </div>
                <div className="level-right">
                    <div className="field is-grouped level-item">
                        <p class="control">
                            <Link to={`${url}/new`} className="button is-link">New Article</Link>
                        </p>
                    </div>
                </div>
            </div>
            <hr />
            {this.props.data &&
            <table className="table is-hoverable is-fullwidth">
                <thead>
                    <tr>
                        <th width="1px">ID</th>
                        <th>Title</th>
                        <th width="1px"></th>
                    </tr>
                </thead>
                <tbody>
                    {this.props.data.map(item => <tr>
                        <td>{item.ID}</td>
                        <td>{item.Title}</td>
                        <td><Link to={`${url}/${item.ID}/edit`} className="button">
                            <i class="fas fa-edit"></i>&nbsp;Edit</Link>
                        </td>
                    </tr>)}
                </tbody>
            </table>
            }
        </div>)
    }
}

export default connect(
    state => ({
        data: state.articles.list,
    }),
    dispatch => bindActionCreators({ fetchList, leaveList }, dispatch),
)(ListPage)
