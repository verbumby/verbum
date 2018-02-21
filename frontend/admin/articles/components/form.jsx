import React from 'react'
import { connect } from 'react-redux'
import { Link } from 'react-router-dom'

class Form extends React.Component {
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

        switch (name) {
            case "DictID":
                value = parseInt(value)
                break;
        }

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
                <label class="label">Dictionary</label>
                <div class="control">
                    <div class="select">
                        <select name="DictID" value={this.state.DictID} onChange={this.handleInputChange} required >
                            <option />
                            {this.props.dicts.map(d => <option value={d.ID}>{d.Title}</option>)}
                        </select>
                    </div>
                </div>
            </div>
            <div className="field">
                <div className="control">
                    <label class="label">Content</label>
                    <textarea
                        className="textarea"
                        type="text"
                        name="content"
                        rows="15"
                        value={this.state.Content}
                        onChange={this.handleInputChange}
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

export default connect(
    state => ({
        dicts: state.config.Dicts
    })
)(Form)
