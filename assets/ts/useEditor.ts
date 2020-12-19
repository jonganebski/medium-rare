import EditorJS, { OutputBlockData } from "@editorjs/editorjs";
import Axios from "axios";
import { EDITORJS_CONFIG } from "./constants";
import { publishBtn, saveStatusEl, unpublishBtn } from "./elements.header";
import { overrideEditorJsStyleBody } from "./page.ReadStory";

let timeoutId: NodeJS.Timeout;
let isCreated = Boolean(document.location.pathname.includes("edit-story"));
let isSaved = false;
let storyId = "";

export const handlePublishBtnClick = async (
  e: Event,
  editor: EditorJS,
  toggle: 1 | -1
) => {
  const publishBtn = e.target as HTMLButtonElement | null;
  if (!publishBtn || !storyId) {
    return;
  }
  publishBtn.disabled = true;
  publishBtn.innerText = "Loading...";
  try {
    clearTimeout(timeoutId);
    !isSaved && (await saveAsDraft(editor));
    const { status } = await Axios.patch(
      `/api/toggle-publish/${storyId}/${toggle}`
    );
    if (status < 300) {
      document.location.href = `/read-story/${storyId}`;
    }
  } catch {
    alert("Failed to publish. Please try again.");
    publishBtn.disabled = false;
    publishBtn.innerText = "Save and Publish";
  }
  return;
};

const saveAsDraft = async (editor: EditorJS) => {
  if (!storyId) {
    return;
  }
  try {
    saveStatusEl && (saveStatusEl.innerText = "Saving...");
    const savedData = await editor.save();
    const { status } = await Axios.patch(`/api/story/${storyId}`, savedData);
    if (status < 300) {
      isSaved = true;
      saveStatusEl && (saveStatusEl.innerText = "Saved");
    }
  } catch {
    isSaved = false;
    saveStatusEl && (saveStatusEl.innerText = "Not Saved");
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
    }
  } catch {
    isCreated = false;
    isSaved = false;
    saveStatusEl && (saveStatusEl.innerText = "Not saved");
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
    }
    const editorEl = document.getElementById(holder);
    const observer = new MutationObserver(() => saveStoryTimeout(editor));
    editorEl &&
      observer.observe(editorEl, {
        attributes: true,
        childList: true,
        subtree: true,
      });
    publishBtn?.addEventListener("click", (e) =>
      handlePublishBtnClick(e, editor, 1)
    );
    unpublishBtn?.addEventListener("click", (e) => {
      handlePublishBtnClick(e, editor, -1);
    });
    window.addEventListener("beforeunload", () => {
      clearTimeout(timeoutId);
      !isSaved && saveAsDraft(editor);
    });
    overrideEditorJsStyleBody();
  });
};
