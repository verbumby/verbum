import React from 'react'

export class InputElement extends React.Component {
    render() {
        return (<div className="field">
            <div className="control">
                {this.props.label ? <label class="label">{this.props.label}</label> : null}
                {this.renderInput({
                    name: this.props.name,
                    value: this.context.formData[this.props.name],
                    onChange: (ev) => this.context.onChange(this.props.name, ev.target.value),
                })}
            </div>
        </div>)
    }

    renderInput({ name, value, onChange }) {
        return <input
            className="input"
            type="text"
            name={name}
            value={value}
            onChange={onChange}
        />
    }
}

InputElement.contextTypes = {
    formData: true,
    onChange: true,
}
