export const BASE_URL =
  process.env.APP_ENV === "DEV"
    ? "http://localhost:4000"
    : process.env.APP_ENV === "PROD"
    ? ""
    : null;
