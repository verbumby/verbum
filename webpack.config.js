const path = require('path');
const webpack = require('webpack');

module.exports = {
  entry: {
    admin: './frontend/admin/index.jsx',
  },
  output: {
    filename: '[name].bundle.js',
    path: path.resolve(__dirname, 'statics')
  },
  optimization: {
    splitChunks: {
      cacheGroups: {
        commons: {
          test: /[\\/]node_modules[\\/]/,
          name: "vendor",
          chunks: "all",
        },
      },
    },
  },
  module: {
    rules: [{
      test: /\.jsx?$/,
      loader: 'babel-loader',
    },{
      test: /\.css$/,
      use: [
        'style-loader',
        'css-loader',
      ],
    }]
  },
  resolve: {
    extensions: ['.js', '.json', '.jsx', '.css']
  },
  plugins: [
    new webpack.HashedModuleIdsPlugin(),
  ],
  performance: {
    maxAssetSize: 10 * 1024 * 1024,
    maxEntrypointSize: 10 * 1024 * 1024,
  },
  target: 'web',
  mode: 'production',
  cache: true,
  devtool: 'source-map',
};
