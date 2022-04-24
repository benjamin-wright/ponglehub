const path = require("path");
const CopyWebpackPlugin = require('copy-webpack-plugin');
const webpack = require('webpack');

module.exports = {
  devServer: {
    hot: true,
    port: 80,
  },
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
    new webpack.DefinePlugin({
      'process.env.NODE_ENV': 'production',
      'process.env.NODE_DEBUG': '"console"'
    })
  ]
};