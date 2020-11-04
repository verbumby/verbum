const path = require('path')

const prod = process.env.NODE_ENV === 'production'
const mode = prod ? 'production' : 'development'
console.log(`BUILD MODE: ${mode}`)

const output = {
    path: path.resolve(__dirname, 'frontend', 'dist'),
    filename: '[name].bundle.js',
}

const server = {
    name: 'server',
    mode,
    entry: {
        server: './frontend/server.ts',
    },
    output,
}

const browser = {
    name: 'browser',
    mode,
    entry: {
        browser: './frontend/browser.ts',
    },
    output,
}

module.exports = [server, browser]
