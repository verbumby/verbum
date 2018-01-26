import React from 'react'
import { Link } from 'react-router-dom'

export * from './form/input-element'

export default class Form extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            formData: props.formData
        }
    }

    getChildContext() {
        return {
            formData: this.state.formData,
            onChange: this.onChildChange.bind(this)
        }
    }

    onChildChange(key, value) {
        this.setState({ formData: { ...this.formData, [key]: value } })
    }

    render() {
        return (<div>
            {this.props.children}
            <div className="field is-grouped">
                <p class="control">
                    <button
                        className="button is-link"
                        onClick={() => this.props.onSave({ formData: this.state.formData })}
                    >
                        Save
                    </button>
                </p>
                <p class="control">
                    <Link className="button" to={this.props.onCancelRedirectTo}>Cancel</Link>
                </p>
            </div>
        </div>)
    }
}

Form.childContextTypes = {
    formData: true,
    onChange: true,
}
