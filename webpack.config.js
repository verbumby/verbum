const path = require('path')
const { CleanWebpackPlugin } = require('clean-webpack-plugin')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const FaviconsWebpackPlugin = require('favicons-webpack-plugin')

const prod = process.env.NODE_ENV === 'production'
const mode = prod ? 'production' : 'development'
console.log(`BUILD MODE: ${mode}`)

const resolve = { extensions: ['.ts', '.tsx', '.js'] }

const server = {
    name: 'server',
    target: 'node15.0',
    mode,
    entry: {
        server: './frontend/server.tsx',
    },
    output: {
        path: path.resolve(__dirname, 'frontend', 'dist'),
        filename: '[name].bundle.js',
    },
    resolve,
    module: {
        rules: [
            {
                test: /\.tsx?$/,
                exclude: /node_modules/,
                use: [
                    {
                        loader: 'ts-loader',
                        options: {
                            compilerOptions: {
                                moduleResolution: 'node',
                            },
                        },
                    },
                ],
            },
        ],
    },
    devtool: 'source-map',
    plugins: [
        new CleanWebpackPlugin({
            cleanOnceBeforeBuildPatterns: ['*.js', '*.map'],
        }),
    ],
}

const browser = {
    name: 'browser',
    target: 'browserslist:> 1%, last 2 versions, Firefox ESR, not dead',
    mode,
    entry: {
        browser: './frontend/browser.tsx',
    },
    output: {
        path: path.resolve(__dirname, 'frontend', 'dist', 'public'),
        filename: '[name].bundle.js',
    },
    resolve,
    module: {
        rules: [
            { test: /\.tsx?$/, use: ['ts-loader'], exclude: /node_modules/ },
        ],
    },
    devtool: 'source-map',
    plugins: [
        new CleanWebpackPlugin(),
        new HtmlWebpackPlugin({
            filename: '../index.html',
            publicPath: '/statics',
            hash: true,
        }),
        new FaviconsWebpackPlugin({
            logo: './frontend/favicon.png',
            prefix: '',
            publicPath: '/statics',
            inject: true,
        }),
    ],
}

module.exports = [server, browser]
