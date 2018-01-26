import React from 'react'

import BaseForm, { InputElement } from '../../core/components/form'

export default class Form extends React.Component {
    render() {
        return (<BaseForm {...this.props}>
            <InputElement label="Dictionary Name" name="title"/>
        </BaseForm>)
    }
}
