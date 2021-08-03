module.exports = {
  publicPath: "/ui",
  configureWebpack: {
    devServer: {
      port: 3000,
      watchOptions: {
        poll: true,
      },
    },
  },
};
