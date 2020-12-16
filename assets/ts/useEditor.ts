import CodeTool from "@editorjs/code";
import EditorJS, {
  LogLevels,
  OutputBlockData,
  OutputData,
} from "@editorjs/editorjs";
import Header from "@editorjs/header";
import ImageTool from "@editorjs/image";
import InlineCode from "@editorjs/inline-code";
import List from "@editorjs/list";
import Quote from "@editorjs/quote";
import Axios from "axios";
import { BASE_URL } from "./constants";
import { publishBtn } from "./elements.header";
import { overrideEditorJsStyleBody } from "./page.ReadStory";

let imgHistory: string[] = [];

// All these imgHistory and requestUnusedPhotosDelete are because editorJS's image plugin does not trigger any event on delete photos!!!!!
// So this app will store all photos list and will request removal of unused photos at the point of publish.
const requestUnusedPhotosDelete = async (
  savedData: OutputData
): Promise<boolean> => {
  const imgBlocks = savedData.blocks.filter((block) => block.type === "image");
  if (imgBlocks.length === 0) {
    return true;
  }
  try {
    imgBlocks.forEach((imgBlock) => {
      const usedImg = imgBlock.data.file.url;
      imgHistory = imgHistory.filter((url) => url !== usedImg);
    });
    const { status: removalStatus } = await Axios.delete("/api/photos", {
      data: { images: Array.from(imgHistory) },
    });
    if (removalStatus < 300) {
      return true;
    }
    return false;
  } catch {
    return false;
  }
};

const handlePublishBtnClick = async (e: Event, editor: EditorJS) => {
  const publishBtn = e.target as HTMLButtonElement | null;
  if (!publishBtn) {
    return;
  }
  publishBtn.disabled = true;
  publishBtn.innerText = "Loading...";
  const savedData = await editor.save();
  if (document.location.pathname.includes("new-story")) {
    try {
      const { status, data: storyId } = await Axios.post(
        "/api/story",
        savedData
      );
      if (status < 300) {
        // request delete images in imgHistory
        if (await !requestUnusedPhotosDelete(savedData)) {
          console.error("There are images failed to remove.");
        }
        document.location.href = `/read-story/${storyId}`;
      }
    } catch {
      alert("Failed to publish. Please try again.");
      publishBtn.disabled = false;
      publishBtn.innerText = "Publish";
    }

    return;
  }
  if (document.location.pathname.includes("edit-story")) {
    const splitedPath = document.location.pathname.split("edit-story");
    const storyId = splitedPath[1].replace(/[/]/g, "");
    try {
      const { status } = await Axios.patch(`/api/story/${storyId}`, savedData);
      if (status < 300) {
        // request delete images in imgHistory
        if (await !requestUnusedPhotosDelete(savedData)) {
          console.error("There are images failed to remove.");
        }
        document.location.href = `/read-story/${storyId}`;
      }
    } catch {
      alert("Failed to update. Please try again.");
      publishBtn.disabled = false;
      publishBtn.innerText = "Save and Publish";
    }
    return;
  }
};

const imageCollector = (holder: string) => {
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
      inlineCode: {
        class: InlineCode,
      },
      quote: Quote,
      list: {
        class: List,
        inlineToolbar: true,
      },
    },
    data: {
      blocks,
    },
    logLevel: LogLevels?.ERROR ?? "ERROR",
  });
  editor.isReady.then(() => {
    publishBtn?.addEventListener("click", (e) =>
      handlePublishBtnClick(e, editor)
    );
    document.body.addEventListener("click", () => imageCollector(holder));
    overrideEditorJsStyleBody();
  });
};
