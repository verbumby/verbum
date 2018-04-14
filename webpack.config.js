const path = require('path');

module.exports = {
  entry: './frontend/admin/index.jsx',
  output: {
    filename: 'admin.js',
    path: path.resolve(__dirname, 'statics')
  },
  module: {
    rules: [{
        test: /\.jsx?$/,
        loader: 'babel-loader',
    }]
  },
  resolve: {
    extensions: ['.js', '.json', '.jsx', '.css']
  },
  target: 'web',
  mode: 'production',
  devtool: 'source-map',
  externals: {
    'react': 'React',
    'react-dom': 'ReactDOM',
    'redux': 'Redux',
    'react-redux': 'ReactRedux',
    'redux-thunk': 'ReduxThunk',
    'react-router-dom': 'ReactRouterDOM',
    'simplemde': 'SimpleMDE',
  }
};
