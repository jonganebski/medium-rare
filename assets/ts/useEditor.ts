import EditorJS, { OutputBlockData, OutputData } from "@editorjs/editorjs";
import Axios from "axios";
import { EDITORJS_CONFIG } from "./constants";
import { publishBtn, saveStatusEl, unpublishBtn } from "./elements/header";
import { overrideEditorJsStyleBody } from "./pages/read-story";

let timeoutId: any;
let isCreated = false;
let isSaved = false;
let storyId = "";
let imgHistory: string[] = [];

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
      imgHistory = [];
      return true;
    }
    return false;
  } catch {
    return false;
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

export const handlePublishBtnClick = async (
  e: Event,
  editor: EditorJS,
  toggle: 1 | -1
) => {
  const publishBtn = e.target as HTMLButtonElement | null;
  if (!publishBtn) {
    return;
  }
  publishBtn.disabled = true;
  publishBtn.innerText = "Loading...";
  try {
    clearTimeout(timeoutId);
    let status;
    if (!isCreated && !storyId) {
      const createdStoryId = await createNewStory(editor);
      ({ status } = await Axios.patch(
        `/api/toggle-publish/${createdStoryId}/${toggle}`
      ));
    } else if (!isSaved) {
      await saveAsDraft(editor);
      ({ status } = await Axios.patch(
        `/api/toggle-publish/${storyId}/${toggle}`
      ));
    }
    if (status < 300) {
      document.location.href = `/read-story/${storyId}`;
    }
  } catch {
    alert("Failed to publish. Please try again.");
    publishBtn.disabled = false;
    publishBtn.innerText = "Publish";
  }
  return;
};

const saveAsDraft = async (editor: EditorJS) => {
  if (!storyId) {
    throw new Error();
  }
  try {
    saveStatusEl && (saveStatusEl.innerText = "Saving...");
    const savedData = await editor.save();
    const { status } = await Axios.patch(`/api/story/${storyId}`, savedData);
    if (status < 300) {
      isSaved = true;
      saveStatusEl && (saveStatusEl.innerText = "Saved");
      requestUnusedPhotosDelete(savedData);
      return;
    }
    throw new Error();
  } catch {
    isSaved = false;
    saveStatusEl && (saveStatusEl.innerText = "Save error");
    throw new Error();
  }
};

const createNewStory = async (editor: EditorJS) => {
  try {
    saveStatusEl && (saveStatusEl.innerText = "Saving...");
    const savedData = await editor.save();
    const { status, data: createdStoryId } = await Axios.post(
      "/api/story",
      savedData
    );
    if (status < 300) {
      storyId = createdStoryId;
      isCreated = true;
      isSaved = true;
      saveStatusEl && (saveStatusEl.innerText = "Saved");
      requestUnusedPhotosDelete(savedData);
      return createdStoryId;
    }
    throw new Error();
  } catch {
    isCreated = false;
    isSaved = false;
    saveStatusEl && (saveStatusEl.innerText = "Not saved");
    throw new Error();
  }
};

const saveStoryTimeout = (editor: EditorJS) => {
  if (!saveStatusEl) {
    return;
  }
  isSaved = false;
  saveStatusEl.innerText = "Not saved";
  clearTimeout(timeoutId);
  timeoutId = setTimeout(() => {
    try {
      if (
        document.location.pathname.includes("new-story") &&
        !isCreated &&
        !storyId
      ) {
        createNewStory(editor);
        return;
      } else {
        saveAsDraft(editor);
      }
    } catch {}
  }, 5000);
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
    if (document.location.pathname.includes("edit-story")) {
      const splitedPath = document.location.pathname.split("edit-story");
      storyId = splitedPath[1].replace(/[/]/g, "");
      isCreated = true;
    }
    const editorEl = document.getElementById(holder);
    const observer = new MutationObserver(() => saveStoryTimeout(editor));
    editorEl &&
      observer.observe(editorEl, {
        attributes: true,
        childList: true,
        subtree: true,
      });
    document.body.addEventListener("click", () => imageCollector(holder));
    publishBtn?.addEventListener("click", (e) =>
      handlePublishBtnClick(e, editor, 1)
    );
    unpublishBtn?.addEventListener("click", (e) => {
      handlePublishBtnClick(e, editor, -1);
    });
    overrideEditorJsStyleBody();
  });
};
