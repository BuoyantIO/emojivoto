const HtmlWebpackPlugin = require("html-webpack-plugin");
const webpack = require('webpack');
const path = require('path');
const dotenv = require('dotenv').config( {
  path: path.join(__dirname, '.env')
} );

module.exports = {
  mode: process.env.NODE_ENV === 'production' ? 'production' : 'development',
  entry: './js/index.js',
  cache: false,
  plugins:
    [
      new HtmlWebpackPlugin({
        template: path.resolve(__dirname, "index.html")
      }),
      new webpack.DefinePlugin( {
        "process.env": dotenv.parsed
      } ),

    ],
  devServer: {
    contentBase: path.join(__dirname, 'dist'),
    watchContentBase: true,
    compress: true,
    disableHostCheck: true,
    port: 8080,
    historyApiFallback: {
      rewrites: [
        {from: /^\/leaderboard/, to: 'index.html'},
      ]
    },
  },
  output: {
    path: path.resolve(__dirname, 'dist'),
    publicPath: 'dist/',
    filename: 'index_bundle.js'
  },
  // devtool: 'inline-cheap-source-map', // uncomment for nicer logging, makes dev slower
  externals: {
    cheerio: 'window',
    'react/addons': 'react',
    'react/lib/ExecutionEnvironment': 'react',
    'react/lib/ReactContext': 'react',
    'react-addons-test-utils': 'react-dom',
  },
  module: {
    rules: [
      { test: /\.js$/, loader: 'babel-loader', exclude: /node_modules/ },
      { test: /\.jsx$/, loader: 'babel-loader', exclude: /node_modules/ },
      {
        test: /\.css$/,
        use: [
          'style-loader',
          { loader: 'css-loader', options: { importLoaders: 1 } },
          'postcss-loader'
        ]
      },
      {
        test: /\.(png|jpg|gif|eot|svg|ttf|woff|woff2)$/,
        use: [
          {
            loader: 'file-loader',
            options: { publicPath: 'dist/' }
          }
        ]
      }
    ]
  }
}
