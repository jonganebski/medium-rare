import EditorJS, {
  LogLevels,
  OutputBlockData,
  OutputData,
} from "@editorjs/editorjs";
import Header from "@editorjs/header";
import CodeTool from "@editorjs/code";
import ImageTool from "@editorjs/image";
import { publishBtn } from "./elements.header";
import Axios from "axios";
import { BASE_URL } from "./constants";

let imgHistory: string[] = [];

// All these imgHistory and requestUnusedPhotosDelete are because editorJS's image plugin does not trigger any event on delete photos!!!!!
// So this app will store all photos list and will request removal of unused photos at the point of publish.
const requestUnusedPhotosDelete = async (
  savedData: OutputData
): Promise<boolean> => {
  const imgBlocks = savedData.blocks.filter((block) => block.type === "image");
  imgBlocks.forEach((imgBlock) => {
    const usedImg = imgBlock.data.file.url;
    imgHistory = imgHistory.filter((url) => url !== usedImg);
  });
  const { status: removalStatus } = await Axios.delete(
    BASE_URL + "/api/photos",
    { data: { images: Array.from(imgHistory) } }
  );
  if (removalStatus === 204) {
    return true;
  }
  return false;
};

const handlePublishBtnClick = async (editor: EditorJS) => {
  const savedData = await editor.save();
  if (document.location.pathname.includes("new-story")) {
    try {
      const { status, data: storyId } = await Axios.post(
        BASE_URL + "/api/story",
        savedData
      );
      if (status === 201) {
        // request delete images in imgHistory
        if (await requestUnusedPhotosDelete(savedData)) {
          document.location.href = `/read-story/${storyId}`;
        }
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
        // request delete images in imgHistory
        if (await requestUnusedPhotosDelete(savedData)) {
          document.location.href = `/read-story/${storyId}`;
        }
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
            byFile: BASE_URL + "/api/photo/byfile",
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
    publishBtn?.addEventListener("click", () => handlePublishBtnClick(editor));
    document.body.addEventListener("click", () => {
      const imgElements = document
        .getElementById(holder)
        ?.querySelectorAll(
          ".image-tool__image-picture"
        ) as NodeListOf<HTMLImageElement>;
      imgElements?.forEach((imgEl) => {
        if (!imgHistory.some((url) => url === imgEl.src)) {
          imgHistory.push(imgEl.src);
        }
      });
    });
  });
};
