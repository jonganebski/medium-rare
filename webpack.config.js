const path = require("path");
const Dotenv = require("dotenv-webpack");
require("dotenv").config();
const MiniCSSExtractPlugin = require("mini-css-extract-plugin");
const webpack = require("webpack");

const ENV = process.env.WEBPACK_ENV;
const ENTRY_FILE = path.resolve(__dirname, "assets", "ts", "main.ts");
const OUTPUT_DIR = path.join(__dirname, "static");

const config = {
  entry: ["@babel/polyfill", ENTRY_FILE],
  mode: ENV,
  plugins: [
    new MiniCSSExtractPlugin({ filename: "styles.css" }),
    new webpack.DefinePlugin({
      "process.env": {
        APP_ENV: JSON.stringify(process.env.APP_ENV),
      },
    }),
    new Dotenv(),
  ],
  module: {
    rules: [
      { test: /\.(js)$/, use: "babel-loader", exclude: "/node_modules" },
      { test: /\.(ts)$/, use: "ts-loader", exclude: "/node_modules" },
      {
        test: /\.(scss)$/,
        use: [
          MiniCSSExtractPlugin.loader,
          "css-loader",
          "postcss-loader",
          "sass-loader",
        ],
      },
    ],
  },
  target: "web",
  resolve: {
    extensions: [".ts", ".tsx", ".js", ".scss"],
    fallback: { fs: false, path: require.resolve("path-browserify") },
  },
  output: { path: OUTPUT_DIR, filename: "[name].js" },
};

module.exports = config;
