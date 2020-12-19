import CodeTool from "@editorjs/code";
import { EditorConfig, LogLevels } from "@editorjs/editorjs";
import Header from "@editorjs/header";
import ImageTool from "@editorjs/image";
import InlineCode from "@editorjs/inline-code";
import List from "@editorjs/list";
import Quote from "@editorjs/quote";

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

export const EDITORJS_CONFIG: EditorConfig = {
  tools: {
    header: {
      class: Header,
      inlineToolbar: true,
      config: {
        levels: [2, 4, 6],
      },
    },
    image: {
      class: ImageTool,
      config: {
        endpoints: {
          byFile: BASE_URL + "/api/photo/byfile",
        },
      },
    },
    code: CodeTool,
    inlineCode: {
      class: InlineCode,
    },
    quote: Quote,
    list: {
      class: List,
      inlineToolbar: true,
    },
  },
  logLevel: LogLevels?.ERROR ?? "ERROR",
};

export const MONTHS = [
  "January",
  "Febuary",
  "March",
  "April",
  "May",
  "June",
  "July",
  "August",
  "September",
  "October",
  "November",
  "December",
];
