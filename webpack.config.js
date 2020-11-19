const path = require('path')
const { CleanWebpackPlugin } = require('clean-webpack-plugin')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const FaviconsWebpackPlugin = require('favicons-webpack-plugin')
const { BundleAnalyzerPlugin } = require('webpack-bundle-analyzer')
const MiniCssExtractPlugin = require('mini-css-extract-plugin')
const CssMinimizerPlugin = require('css-minimizer-webpack-plugin')

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
                use: [{ loader: 'ts-loader' }],
            },
            { test: /\.css$/, use: 'null-loader'},
        ],
    },
    devtool: 'source-map',
    plugins: [
        new CleanWebpackPlugin({
            cleanOnceBeforeBuildPatterns: ['*.js', '*.map'],
        }),
        new BundleAnalyzerPlugin({
            analyzerMode: 'static',
            reportFilename: 'server_report.html',
            openAnalyzer: false,
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
        filename: '[name].[contenthash].bundle.js',
    },
    resolve,
    module: {
        rules: [
            { test: /\.tsx?$/, use: ['ts-loader'], exclude: /node_modules/ },
            { test: /\.css$/, use: [MiniCssExtractPlugin.loader, 'css-loader']},
        ],
    },
    devtool: 'source-map',
    plugins: [
        new CleanWebpackPlugin(),
        new HtmlWebpackPlugin({
            template: './frontend/index.html',
            filename: '../index.html',
            publicPath: '/statics',
        }),
        new FaviconsWebpackPlugin({
            logo: './frontend/favicon.png',
            prefix: prod ? 'favicon-[contenthash]' : 'favicon',
            publicPath: '/statics/',
            inject: true,
        }),
        new BundleAnalyzerPlugin({
            analyzerMode: 'static',
            reportFilename: '../browser_report.html',
            openAnalyzer: false,
        }),
        new MiniCssExtractPlugin({
            filename: '[name].[contenthash].bundle.css'
        }),
    ],
    optimization: {
        minimizer: [
            `...`,
            new CssMinimizerPlugin(),
        ],
    },
}

module.exports = [server, browser]
