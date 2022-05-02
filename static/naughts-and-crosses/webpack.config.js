const path = require("path");
const CopyWebpackPlugin = require('copy-webpack-plugin');
const NodePolyfillPlugin = require("node-polyfill-webpack-plugin");
const webpack = require('webpack');

module.exports = {
  target: "web",
  entry: {
    index: "./src/views/index.ts"
  },
  output: {
    path: path.resolve(__dirname, "dist"),
    filename: "js/[name].bundle.js",
  },
  module: {
    rules: [
      {
        test: /\.tsx?$/,
        use: 'ts-loader',
        exclude: /node_modules/,
      },
      {
        test: /\.css$/i,
        use: ["style-loader", "css-loader"],
      }
    ]
  },
  resolve: {
    extensions: ['.tsx', '.ts', '.js', '.css'],
  },
  plugins: [
    new CopyWebpackPlugin({
        patterns: [
            { from: './static', to: './' }
        ]
    }),
    new NodePolyfillPlugin(),
    new webpack.DefinePlugin({
      'process.env.NODE_DEBUG': '"console"'
    })
  ]
};