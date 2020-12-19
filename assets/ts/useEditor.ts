import EditorJS, { OutputBlockData } from "@editorjs/editorjs";
import Axios from "axios";
import { EDITORJS_CONFIG } from "./constants";
import { publishBtn, saveStatusEl, unpublishBtn } from "./elements.header";
import { overrideEditorJsStyleBody } from "./page.ReadStory";

let timeoutId: NodeJS.Timeout;
let isCreated = false;
let isSaved = false;
let storyId = "";

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
    publishBtn?.addEventListener("click", (e) =>
      handlePublishBtnClick(e, editor, 1)
    );
    unpublishBtn?.addEventListener("click", (e) => {
      handlePublishBtnClick(e, editor, -1);
    });
    // window.addEventListener("beforeunload", () => {
    //   clearTimeout(timeoutId);
    //   !isSaved && saveAsDraft(editor);
    // });
    overrideEditorJsStyleBody();
  });
};
