module.exports = {
  plugins: [
      'gatsby-theme-docz',
      `gatsby-plugin-sharp`,
      {
          resolve: `gatsby-source-filesystem`,
          options: {
              path: `${__dirname}/src/pages`,
          },
      },
  ]
}