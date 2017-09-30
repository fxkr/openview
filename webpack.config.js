const path = require('path');
const webpack = require('webpack');

module.exports = {

  context: path.resolve(__dirname, './frontend'),

  entry: {
    app: './app.js',
  },

  output: {
    path: path.resolve(__dirname, './dist'),
    filename: 'openview.bundle.js',
    publicPath: '/',
  },

  devServer: {
    contentBase: path.resolve(__dirname, './frontend'),
    proxy: {
      '/images': {
        target: 'http://localhost:3000/',
        secure: false
      },
      '/api': {
        target: 'http://localhost:3000/',
        secure: false
      },
    },
  },

  module: {
    loaders: [
      {test: /\.html/, loader: "file-loader?name=[name].[ext]"},
      {test: /\.css$/, loader: "style-loader!css-loader"},
      {test: /\.png$/, loader: "url-loader?limit=100000"},
      {test: /\.jpg$/, loader: "file-loader"},
      {test: /\.gif/, loader: "file-loader"},
      {
        test: /\.(woff|woff2)(\?v=\d+\.\d+\.\d+)?$/,
        loader: 'url-loader?limit=10000&mimetype=application/font-woff'
      },
      {test: /\.ttf(\?v=\d+\.\d+\.\d+)?$/, loader: 'url-loader?limit=10000&mimetype=application/octet-stream'},
      {test: /\.eot(\?v=\d+\.\d+\.\d+)?$/, loader: 'file'},
      {test: /\.svg(\?v=\d+\.\d+\.\d+)?$/, loader: 'url-loader?limit=10000&mimetype=image/svg+xml'},
      {
        test: /\.styl$/,
        loader: 'style-loader!css-loader!stylus-loader?paths=node_modules/bootstrap-stylus/stylus/'
      },
      {test: /.js?$/, loader: 'babel-loader', exclude: /node_modules/, query: {presets: ['env']}},
      {test: /.js?$/, loader: 'eslint-loader', exclude: /node_modules/},
    ]
  },

  plugins: [
    new webpack.ProvidePlugin({ // Also see .eslintrc
      $: "jquery",
      jQuery: "jquery",

      PhotoSwipe: 'photoswipe',
      PhotoSwipeUI_Default: 'photoswipe/src/js/ui/photoswipe-ui-default.js'
    })
  ],

  cache: true,
};