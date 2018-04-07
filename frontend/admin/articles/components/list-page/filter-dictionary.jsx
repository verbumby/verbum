import React from 'react'
import { connect } from 'react-redux'

class FilterDictionary extends React.Component {
    render() {
        const { value, dicts, onChange } = this.props
        return (<div class="select">
            <select value={value} onChange={(ev) => onChange(parseInt(ev.target.value))}>
                <option value="-1">- Filter by Dictionary -</option>
                {dicts.map(d => <option value={d.ID}>{d.Title}</option>)}
            </select>
        </div>)
    }
}

export default connect(
    state => ({
        dicts: state.config.Dicts
    })
)(FilterDictionary)
