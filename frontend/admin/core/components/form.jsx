import React from 'react'
import { Link } from 'react-router-dom'

export * from './form/input-element'

export default class Form extends React.Component {
    render() {
        return (<div>
            {this.props.children}
            <div className="field is-grouped">
                <p class="control">
                    <button className="button is-link">Save</button>
                </p>
                <p class="control">
                    <Link className="button" to={this.props.onCancelRedirectTo}>Cancel</Link>
                </p>
            </div>
        </div>)
    }
}
