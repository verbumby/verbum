import React from 'react'
import { connect } from 'react-redux'
import { Link, Route, withRouter } from 'react-router-dom'

class TopNav extends React.Component {
    render() {
        return (<nav className="navbar is-link">
            <div className="navbar-brand">
                <Link className="navbar-item has-text-weight-bold is-size-4" to="/">Verbum Admin</Link>
            </div>
            <div className="navbar-menu">
                <div className="navbar-start">
                    <NavLink path="/dictionaries">Dictionaries</NavLink>
                    <NavLink path="/articles">Articles</NavLink>
                </div>
                <div className="navbar-end">
                    <div className="navbar-item">
                        {this.props.user.name}
                    </div>
                </div>
            </div>
        </nav>)
    }
}

const NavLink = ({ path, children }) => <Route path={path} children={({ match }) => (
    <Link className={`navbar-item ${match ? 'is-active' : ''}`} to={path}>{children}</Link>
)} />

export default withRouter(connect(
    store => ({
        user: store.user
    }),
)(TopNav))
