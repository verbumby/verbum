import React from 'react'

export class InputElement extends React.Component {
    render() {
        return (<div className="field">
            <div className="control">
                {this.props.label ? <label class="label">{this.props.label}</label> : null}
                <input className="input" type="text" />
            </div>
        </div>)
    }
}
