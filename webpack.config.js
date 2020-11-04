const path = require('path')

const prod = process.env.NODE_ENV === 'production'
const mode = prod ? 'production' : 'development'
console.log(`BUILD MODE: ${mode}`)

const output = {
    path: path.resolve(__dirname, 'frontend', 'dist'),
    filename: '[name].bundle.js',
}

const resolve = { extensions: ['.ts', '.tsx', '.js'] }

const server = {
    name: 'server',
    mode,
    entry: {
        server: './frontend/server.tsx',
    },
    output,
    resolve,
    module: {
        rules: [
            { test: /\.tsx?$/, use: ['ts-loader'], exclude: /node_modules/ },
        ],
    },
}

const browser = {
    name: 'browser',
    mode,
    entry: {
        browser: './frontend/browser.tsx',
    },
    output,
    resolve,
    module: {
        rules: [
            { test: /\.tsx?$/, use: ['ts-loader'], exclude: /node_modules/ },
        ],
    },
}

module.exports = [server, browser]
