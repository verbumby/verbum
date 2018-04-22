import React from 'react'
import { connect } from 'react-redux'
import { Link, Route, withRouter } from 'react-router-dom'

class TopNav extends React.Component {
    render() {
        return (<nav className="navbar navbar-expand-sm navbar-dark" style={{backgroundColor:'#28a745f0'}}>
            <Link className="navbar-brand mb-0 h1" to="/">Verbum</Link>
            <button className="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarSupportedContent" aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation">
                <span className="navbar-toggler-icon"></span>
            </button>

            <div className="collapse navbar-collapse" id="navbarSupportedContent">
                <ul className="navbar-nav mr-auto">
                    <NavLink path="/dictionaries">Dictionaries</NavLink>
                    <NavLink path="/articles">Articles</NavLink>
                </ul>
                <span className="navbar-text">{this.props.user.name}</span>
            </div>
        </nav>)
    }
}

const NavLink = ({ path, children }) => <Route path={path} children={({ match }) => (
    <li className={`nav-item ${match ? 'active' : ''}`}>
        <Link className="nav-link" to={path}>{children}</Link>
    </li>
)} />

export default withRouter(connect(
    store => ({
        user: store.user
    }),
)(TopNav))
