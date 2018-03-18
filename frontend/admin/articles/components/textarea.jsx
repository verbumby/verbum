import React from 'react'
import SimpleMDE from 'simplemde'

export default class Textarea extends React.Component {
    shouldComponentUpdate(nextProps, nextState) {
        return nextProps.value != this.simplemde.value()
    }

    setElement(element) {
        this.element = element
        if (element) {
            this.simplemde = this.newSimpleMDE(element)
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

    newSimpleMDE(element) {
        const smde = new SimpleMDE({
            element,
            autofocus: true,
            spellChecker: false,
            toolbar: [
                {
                    name: "bold",
                    action: SimpleMDE.toggleBold,
                    title: "Bold",
                    className: "fa fa-bold",
                },
                {
                    name: "italic",
                    action: SimpleMDE.toggleItalic,
                    title: "Italic",
                    className: "fa fa-italic",
                },
                {
                    name: "strikethrough",
                    action: SimpleMDE.toggleStrikethrough,
                    title: "Strikethrough",
                    className: "fa fa-strikethrough",
                },
                {
                    name: "quote",
                    action: SimpleMDE.toggleBlockquote,
                    title: "Quote",
                    className: "fa fa-quote-left",
                },
                "|",
                {
                    name: "unordered-list",
                    action: SimpleMDE.toggleUnorderedList,
                    title: "Generic List",
                    className: "fa fa-list-ul",
                },
                {
                    name: "ordered-list",
                    action: SimpleMDE.toggleOrderedList,
                    title: "Numbered List",
                    className: "fa fa-list-ol",
                },
                "|",
                {
                    name: "headword",
                    action: this.headwordAction,
                    title: "Headword",
                    className: "fa fa-header color-tag",
                }
            ],
        })
        smde.codemirror.addKeyMap({
            'Shift-Alt-W': this.headwordAction,
        }, true)
        return smde
    }

    headwordAction(codemirror) {
        const selection = codemirror.getSelection()
        codemirror.replaceSelection(`<v-hw>${selection}</v-hw>`, 'around')
    }
}
