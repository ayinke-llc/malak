// const config = {
//   plugins: {
//     tailwindcss: {},
//   },
// };
//
// export default config;
//

/** @type {import('postcss-load-config').Config} */
export const config = {
  plugins: {
    "postcss-import": {},
    "tailwindcss/nesting": {},
    tailwindcss: {},
    autoprefixer: {},
  },
};
