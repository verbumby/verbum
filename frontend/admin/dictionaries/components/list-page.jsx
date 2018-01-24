import React from 'react'
import { Link } from 'react-router-dom'

export default class ListPage extends React.Component {
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


            Dict list
        </div>)
    }
}
