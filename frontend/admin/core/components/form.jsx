import React from 'react'

export default class Form extends React.Component {
    render() {
        return (<div>
            {this.props.children}
            <div className="field is-grouped">
                <p class="control">
                    <button className="button is-link">Save</button>
                </p>
                <p class="control">
                    <button class="button">Cancel</button>
                </p>
            </div>
        </div>)
    }
}
