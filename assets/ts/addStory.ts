import EditorJS from "@editorjs/editorjs";
import Header from "@editorjs/header";
import CodeTool from "@editorjs/code";
import ImageTool from "@editorjs/image";
import { publishBtn } from "./elements.header";
import Axios from "axios";
import { BASE_URL } from "./constants";

const initEditor = () => {
  const editor = new EditorJS({
    holder: "editor__container",
    placeholder: "Write your story",
    tools: {
      header: {
        class: Header,
        inlineToolbar: true,
        config: {
          levels: [2, 4, 6],
        },
      },
      code: CodeTool,
      image: {
        class: ImageTool,
        config: {
          endpoints: {
            byFile: "http://localhost:4000/upload/photo/byfile",
          },
        },
      },
    },
    data: {
      blocks: [
        {
          type: "header",
          data: { level: 2, text: "Almost before we knew it" },
        },
        {
          type: "header",
          data: { level: 4, text: "Almost before we knew it" },
        },
      ],
    },
  });
  publishBtn?.addEventListener("click", async () => {
    const savedData = await editor.save();
    console.log(savedData);
    const response = await Axios.post(BASE_URL + "/upload/story", savedData);
    console.log(response);
  });
};

const addStory = () => {
  if (document.location.pathname.includes("new-story")) {
    initEditor();
  }
};

addStory();
