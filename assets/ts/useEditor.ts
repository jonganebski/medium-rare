import EditorJS, { OutputBlockData, OutputData } from "@editorjs/editorjs";
import Axios from "axios";
import { EDITORJS_CONFIG } from "./constants";
import { publishBtn } from "./elements.header";
import { overrideEditorJsStyleBody } from "./page.ReadStory";

let imgHistory: string[] = [];
let isPublishBtn = false;
// All these imgHistory and requestUnusedPhotosDelete etc are because editorJS's image plugin does not trigger any event on delete photos!!!!!
// So this app will store all photos list and will request removal of unused photos at the point of publish or beforeunload event.
const requestUnusedPhotosDelete = async (
  savedData: OutputData
): Promise<boolean> => {
  const imgBlocks = savedData.blocks.filter((block) => block.type === "image");
  try {
    imgBlocks.forEach((imgBlock) => {
      const usedImg = imgBlock.data.file.url;
      imgHistory = imgHistory.filter((url) => url !== usedImg);
    });
    if (imgHistory.length === 0) {
      return true;
    }
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

const requestAllPhotosDelete = (e: Event) => {
  // this happens beforeunload. so I cannot do anything even if it fails :(
  e.preventDefault();
  if (!isPublishBtn) {
    Axios.delete("/api/photos", {
      data: { images: Array.from(imgHistory) },
    });
  }
};

const handlePublishBtnClick = async (e: Event, editor: EditorJS) => {
  isPublishBtn = true;
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
      isPublishBtn = false;
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
    ...EDITORJS_CONFIG,
    holder,
    placeholder,
    data: {
      blocks,
    },
  });
  editor.isReady.then(() => {
    publishBtn?.addEventListener("click", (e) =>
      handlePublishBtnClick(e, editor)
    );
    document.body.addEventListener("click", () => imageCollector(holder));
    window.addEventListener("beforeunload", requestAllPhotosDelete, false);
    overrideEditorJsStyleBody();
  });
};
