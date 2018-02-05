import React from 'react'

import { InputElement } from './input-element'

export class TextareaElement extends InputElement {
    renderInput({ name, value, onChange }) {
        return <textarea
            className="textarea"
            type="text"
            name={name}
            value={value}
            onChange={onChange}
        />
    }
}
