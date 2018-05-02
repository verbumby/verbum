import React from 'react'
import { connect } from 'react-redux'
import { Link } from 'react-router-dom'

import './form.css'
import Textarea from './textarea'
import Task from './form/task'

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
            Status: task.Status == 'PENDING' ? 'DONE' : 'PENDING',
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
            <div class="form-group">
                <label class="label">Tasks</label>
                {this.state.Tasks.map((task, index) => {
                    return <Task onToggle={() => { this.toggleTask(index) }} task={task} index={index + 1} />
                })}
            </div>
            <div class="form-group">
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

            <div className="form-row">
                <button
                    className="btn btn-primary"
                    type="submit"
                    ref={button => this.submitButton = button}
                >
                    Save
                </button>
            </div>
        </form>
    }
}

export default connect(
    state => ({
        dicts: state.config.Dicts
    })
)(Form)
