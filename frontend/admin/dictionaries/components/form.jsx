import React from 'react'
import { Link } from 'react-router-dom'

export default class Form extends React.Component {
    constructor(props) {
        super(props)
        this.state = props.formData

        this.handleInputChange = this.handleInputChange.bind(this)
        this.handleSubmit = this.handleSubmit.bind(this)
    }

    handleInputChange(event) {
        const target = event.target
        let value = target.type === 'checkbox' ? target.checked : target.value
        const name = target.name

        this.setState({
            [name]: value
        })
    }

    handleSubmit(event) {
        event.preventDefault()
        this.props.onSave({ formData: this.state })
    }

    render() {
        return <form onSubmit={this.handleSubmit}>
            <div class="field">
                <label class="label">Title</label>
                <div class="control">
                    <input className="input"
                        type="text"
                        name="Title"
                        value={this.state.Title}
                        onChange={this.handleInputChange}
                        required
                    />
                </div>
            </div>
            <div className="field is-grouped">
                <p class="control">
                    <button className="button is-link" type="submit">Save</button>
                </p>
                <p class="control">
                    <Link className="button" to={this.props.onCancelRedirectTo}>Cancel</Link>
                </p>
            </div>
        </form>
    }
}
