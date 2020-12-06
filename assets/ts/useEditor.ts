import EditorJS, { OutputBlockData } from "@editorjs/editorjs";
import Header from "@editorjs/header";
import CodeTool from "@editorjs/code";
import ImageTool from "@editorjs/image";
import { publishBtn } from "./elements.header";
import Axios from "axios";
import { BASE_URL } from "./constants";

export const useEditor = (
  holder: string,
  placeholder: string,
  blocks: OutputBlockData[]
) => {
  const editor = new EditorJS({
    holder,
    placeholder,
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
            byFile: "http://localhost:4000/api/photo/byfile",
          },
        },
      },
    },
    data: {
      blocks,
    },
  });
  publishBtn?.addEventListener("click", async () => {
    const savedData = await editor.save();
    console.log(savedData);
    if (document.location.pathname.includes("read")) {
      const response = await Axios.post(BASE_URL + "/api/story", savedData);
      console.log(response);
      return;
    }
    if (document.location.pathname.includes("edit-story")) {
      const split = document.location.pathname.split("/");
      const storyId = split[split.length - 1];
      const response = await Axios.patch(
        BASE_URL + `/api/story/${storyId}`,
        savedData
      );
      console.log(response);
      return;
    }
  });
};
