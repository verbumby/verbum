import React from 'react'
import { connect } from 'react-redux'

class Messages extends React.Component {
    render() {
        return (<div className="messages">
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
