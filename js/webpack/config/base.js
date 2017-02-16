if (!global.Promise) {
    global.Promise = require('bluebird');
}

var webpack = require('webpack');
var StyleLintPlugin = require('stylelint-webpack-plugin');
var path = require('path');
var base = path.join(__dirname, '..', '..', '..');
var nodeEnv = process.env.NODE_ENV === 'production' ? 'production' : 'development';

module.exports = function() {
    return {
        output: {
            path: path.join(base, 'js', 'public'),
            publicPath: '/assets/',
            filename: '[name].js',
        },

        module: {
            rules: [
                {
                    test: /\.js?$/,
                    enforce: 'pre',
                    exclude: /node_modules/,
                    use: ['eslint-loader'],
                },
                {
                    test: /\.js?$/,
                    exclude: /node_modules/,
                    use: ['babel-loader'],
                },
                {
                    test: /\.css$/,
                    use: [
                        'style-loader/useable',
                        { loader: 'css-loader', options: { importLoaders: 1 } },
                        {
                            loader: 'postcss-loader',
                            options: {
                                ident: 'postcss',
                                plugins: function() {
                                    return [
                                        require('postcss-nested'),
                                        require('autoprefixer'),
                                    ];
                                },
                            },
                        },
                    ],
                },
                {
                    test: /\.woff($|\?)/,
                    use: [{ loader: 'url-loader', options: { limit: 10000 } }],
                },
                {
                    test: /\.woff2($|\?)/,
                    use: [{ loader: 'url-loader', options: { limit: 10000 } }],
                },
                {
                    test: /\.otf($|\?)/,
                    use: [{ loader: 'url-loader', options: { limit: 10000 } }],
                },
                {
                    test: /\.ttf($|\?)/,
                    use: [{ loader: 'url-loader', options: { limit: 10000 } }],
                },
                {
                    test: /\.eot($|\?)/,
                    use: [{ loader: 'url-loader', options: { limit: 10000 } }],
                },
                {
                    test: /\.svg($|\?)/,
                    use: [{ loader: 'url-loader', options: { limit: 10000 } }],
                },
                {
                    test: /\.png$/,
                    use: [{ loader: 'url-loader', options: { limit: 10000 } }],
                },
                {
                    test: /\.jpg$/,
                    use: [{ loader: 'url-loader', options: { limit: 10000 } }],
                },
                {
                    test: /\.gif$/,
                    use: [{ loader: 'url-loader', options: { limit: 10000 } }],
                },
            ],
        },

        resolve: {
            extensions: ['.webpack.js', '.web.js', '.js'],
            alias: {
                js: path.join(base, 'js'),
            },
        },

        plugins: [
            new webpack.DefinePlugin({
                'process.env.NODE_ENV': JSON.stringify(nodeEnv),
            }),
            new StyleLintPlugin({
                configFile: '.stylelintrc',
                files: ['./js/**/*.css'],
                failOnError: false,
            }),
        ],
    };
};
