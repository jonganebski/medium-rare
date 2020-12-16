export const BASE_URL =
  process.env.APP_ENV === "DEV"
    ? "http://localhost:4000"
    : process.env.APP_ENV === "PROD"
    ? "https://go-medium-rare.herokuapp.com"
    : null;

export const INITIAL_BLOCKS = [
  {
    type: "header",
    data: { level: 2, text: "Title" },
  },
  {
    type: "paragraph",
    data: { text: "Write your story" },
  },
];
