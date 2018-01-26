import React from 'react'

import BaseForm, { InputElement } from '../../core/components/form'

export default class Form extends React.Component {
    render() {
        return (<BaseForm {...this.props} formData={{}} onSave={({formData}) => console.log(formData)}>
            <InputElement label="Dictionary Name" name="title"/>
        </BaseForm>)
    }
}
