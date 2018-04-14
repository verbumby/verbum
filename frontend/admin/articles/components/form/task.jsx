import React from 'react'

export default class Task extends React.Component {
    constructor(props) {
        super(props)
        this.keyUpHandler = this.keyUpHandler.bind(this)
    }

    componentWillMount() {
        document.addEventListener('keyup', this.keyUpHandler)
    }

    componentWillUnmount() {
        document.removeEventListener('keyup', this.keyUpHandler)
    }

    keyUpHandler(e) {
        const { index, onToggle } = this.props
        if (e.ctrlKey && e.keyCode - 48 == index ) {
            e.preventDefault()
            e.stopPropagation()
            onToggle()
        }
    }

    render() {
        const { onToggle, task: it, index: i } = this.props

        const style = it.Status == 'PENDING' ? 'is-info' : 'is-success'
        return (<div>
            <a class={`button ${style}`} onClick={onToggle}>
                ({i})&nbsp;{it.Task.Title}
            </a>
        </div>)
    }
}
