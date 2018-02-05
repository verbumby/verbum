import React from 'react'

import BaseForm, { TextareaElement } from '../../core/components/form'

export default class Form extends React.Component {
    render() {
        return (<BaseForm {...this.props}>
            <TextareaElement label="Content" name="content"/>
        </BaseForm>)
    }
}
