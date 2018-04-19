import React from 'react'
import { connect } from 'react-redux'

import './messages.css'

class Messages extends React.Component {
    render() {
        return (<div className="messages position-absolute w-100">
            {this.props.messages.map(message => (
                <div key={message.id} className={`alert alert-${message.level}`}>
                    {message.message}
                </div>
            ))}
        </div>)
    }
}

export default connect(
    state => ({
        messages: state.messages
    })
)(Messages)
