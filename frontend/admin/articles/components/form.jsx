import React from 'react'
import { connect } from 'react-redux'
import { Link } from 'react-router-dom'

import Textarea from './textarea'

class Form extends React.Component {
    constructor(props) {
        super(props)
        this.state = props.formData

        this.handleInputChange = this.handleInputChange.bind(this)
        this.handleSubmit = this.handleSubmit.bind(this)
        this.toggleTask = this.toggleTask.bind(this)
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

    toggleTask(n) {
        let task = this.state.Tasks[n]
        task = {
            ...task,
            Status: task.Status == 'PENDING' ? 'DONE': 'PENDING',
        }

        const Tasks = [...this.state.Tasks]
        Tasks[n] = task

        this.setState({ Tasks })
    }

    handleSubmit(event) {
        event.preventDefault()
        this.props.onSave({ formData: this.state })
    }

    render() {
        return <form onSubmit={this.handleSubmit}>
            <div class="columns">
                <div class="column">
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
                </div>
                <div class="column">
                    <div class="field">
                        <label class="label">Tasks</label>
                        {this.state.Tasks.map((it, i) => {
                            const style = it.Status == 'PENDING' ? 'is-info' : 'is-success'
                            const icon = it.Status == 'PENDING'
                                ? <i class="fa fa-circle-o" aria-hidden="true"></i>
                                : <i class="fa fa-check-circle" aria-hidden="true"></i>

                            return <div>
                                <a class={`button ${style}`} onClick={() => {this.toggleTask(i)}}>{icon}&nbsp;{it.Task.Title}</a>
                            </div>
                        })}
                    </div>
                </div>
            </div>
            <div className="field">
                <div className="control">
                    <label class="label">Content</label>
                    <Textarea
                        className="textarea"
                        type="text"
                        name="Content"
                        rows="15"
                        value={this.state.Content}
                        onChange={this.handleInputChange}
                        onSave={() => { this.submitButton.click() }}
                    />
                </div>
            </div>
            <div className="field is-grouped">
                <p class="control">
                    <button className="button is-link" type="submit"
                        ref={button => this.submitButton = button}>
                        Save
                    </button>
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
