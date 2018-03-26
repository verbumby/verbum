import React from 'react'
import { Link } from 'react-router-dom'
import { connect } from 'react-redux'
import { bindActionCreators } from 'redux'

import { fetchList, leaveList } from '../actions'
import { parseURLSearchParams, assembleURLQuery } from '../../utils'

class ListPage extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            offset: 0,
            limit: 30,
        }
    }

    componentWillMount() {
        this.fetchList()
    }

    componentWillUnmount() {
        this.props.leaveList()
    }

    fetchList() {
        this.props.fetchList({
            offset: this.state.offset,
            limit: this.state.limit,
        })
    }

    render() {
        const { url } = this.props.match
        const onPrevPageClick = () => {
            const offset = Math.max(0, this.state.offset - this.state.limit)
            this.setState({ offset }, () => { this.fetchList() })
        }

        const onNextPageClick = () => {
            const offset = this.state.offset + this.state.limit
            this.setState({ offset }, () => { this.fetchList() })
        }

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
            <div>
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
                        <td><Link to={`${url}/${item.ID}/edit`} className="button">Edit</Link>
                        </td>
                    </tr>)}
                </tbody>
            </table>
            <nav class="pagination" role="navigation" aria-label="pagination">
                <a class="pagination-previous" onClick={onPrevPageClick}>Previous</a>
                <a class="pagination-next" onClick={onNextPageClick}>Next page</a>
            </nav>
            </div>
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
