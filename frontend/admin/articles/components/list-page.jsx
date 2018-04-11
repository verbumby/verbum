import React from 'react'
import { Link } from 'react-router-dom'
import { connect } from 'react-redux'
import { bindActionCreators } from 'redux'

import FilterDictionary from './list-page/filter-dictionary'
import FilterTask from './list-page/filter-task'
import { fetchList, leaveList } from '../actions'
import { parseURLSearchParams, assembleURLQuery } from '../../utils'

class ListPage extends React.Component {
    constructor(props) {
        super(props)

        this.state = {
            ...this.getDefaultState(),
            ...JSON.parse(localStorage.getItem('articles.list-page-state') || '{}')
        }
    }

    setState(props, callback) {
        return super.setState(props, () => {
            localStorage.setItem('articles.list-page-state', JSON.stringify(this.state))
            callback()
        })
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
            filter$DictID: this.state.filter$DictID,
            filter$TitlePrefix: this.state.filter$TitlePrefix,
            filter$TaskID: this.state.filter$TaskID,
            _defaults: this.getDefaultState(),
        })
    }

    getDefaultState() {
        return {
            offset: 0,
            limit: 20,
            filter$DictID: -1,
            filter$TitlePrefix: '',
            filter$TaskID: -1,
        }
    }

    setFilterState(state) {
        this.setState({ ...state, offset: 0 }, () => { this.fetchList() })
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
            {/* filter */}
            <div class="field is-grouped">
                <p class="control">
                    <input
                        class="input"
                        type="text"
                        value={this.state.filter$TitlePrefix}
                        onChange={ev => this.setFilterState({ filter$TitlePrefix: ev.target.value })}
                    />
                </p>
                <p class="control">
                    <FilterDictionary value={this.state.filter$DictID} onChange={filter$DictID => this.setFilterState({ filter$DictID })} />
                </p>
                <p class="control">
                    <FilterTask value={this.state.filter$TaskID} onChange={filter$TaskID => this.setFilterState({ filter$TaskID })} />
                </p>
            </div>

            {this.props.data &&
                <div>
                    {/* table */}
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
                </div>
            }
            <nav class="pagination" role="navigation" aria-label="pagination">
                <a class="pagination-previous" onClick={onPrevPageClick}>Previous</a>
                <a class="pagination-next" onClick={onNextPageClick}>Next page</a>
            </nav>
        </div>)
    }
}

export default connect(
    state => ({
        data: state.articles.list,
    }),
    dispatch => bindActionCreators({ fetchList, leaveList }, dispatch),
)(ListPage)
