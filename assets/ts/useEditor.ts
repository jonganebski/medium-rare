import EditorJS, { LogLevels, OutputBlockData } from "@editorjs/editorjs";
import Header from "@editorjs/header";
import CodeTool from "@editorjs/code";
import ImageTool from "@editorjs/image";
import { publishBtn } from "./elements.header";
import Axios from "axios";
import { BASE_URL } from "./constants";

const handlePublishBtnClick = async (
  editor: EditorJS,
  imgUrlsHistory: Set<string>
) => {
  const savedData = await editor.save();
  const imgBlocks = savedData.blocks.filter((block) => block.type === "image");
  imgBlocks.forEach((imgBlock) => {
    const usedImg = imgBlock.data.file.url;
    imgUrlsHistory.delete(usedImg);
  });
  if (document.location.pathname.includes("new-story")) {
    try {
      const { status, data: storyId } = await Axios.post(
        BASE_URL + "/api/story",
        savedData
      );
      if (status === 201) {
        document.location.href = `/read/${storyId}`;
        // request delete images in imgUrlsHistory
      }
    } catch {}
    return;
  }
  if (document.location.pathname.includes("edit-story")) {
    const splitedPath = document.location.pathname.split("edit-story");
    const storyId = splitedPath[1].replace(/[/]/g, "");
    try {
      const { status } = await Axios.patch(
        BASE_URL + `/api/story/${storyId}`,
        savedData
      );
      if (status === 200) {
        document.location.href = `/read/${storyId}`;
        // request delete images in imgUrlsHistory
      }
    } catch {}
    return;
  }
};

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
      image: {
        class: ImageTool,
        config: {
          endpoints: {
            byFile: "http://localhost:4000/api/photo/byfile",
          },
        },
      },
      code: CodeTool,
    },
    data: {
      blocks,
    },
    logLevel: LogLevels?.ERROR ?? "ERROR",
  });
  editor.isReady.then(() => {
    const imgUrlsHistory = new Set<string>();
    publishBtn?.addEventListener("click", () =>
      handlePublishBtnClick(editor, imgUrlsHistory)
    );
    document.body.addEventListener("click", () => {
      const imgElements = document
        .getElementById(holder)
        ?.querySelectorAll(
          ".image-tool__image-picture"
        ) as NodeListOf<HTMLImageElement>;
      imgElements?.forEach((imgEl) => {
        imgUrlsHistory.add(imgEl.src);
      });
    });
  });
};
