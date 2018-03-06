import React from 'react'
import SimpleMDE from 'simplemde'

export default class Textarea extends React.Component {
    shouldComponentUpdate(nextProps, nextState) {
        return nextProps.value != this.simplemde.value()
    }

    setElement(element) {
        this.element = element
        if (element) {
            this.simplemde = new SimpleMDE({
                element,
                autofocus: true,
                spellChecker: false,
            })
            this.simplemde.codemirror.on('change', () => {
                this.props.onChange({
                    target: {
                        name: this.props.name,
                        value: this.simplemde.value(),
                    }
                })
            })
        } else {
            this.simplemde.toTextArea()
            this.simplemde = null
        }
    }
    render() {
        return <textarea {...this.props} ref={(el) => {this.setElement(el)}}/>
    }
}
