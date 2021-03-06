import EditorJS, { OutputBlockData } from "@editorjs/editorjs";
import Axios from "axios";
import { EDITORJS_CONFIG } from "../constants";
import { deleteComment } from "../deleteComment";
import { deleteStory } from "../deleteStory";
import {
  bookmarkContainer,
  commentDrawer,
  commentDrawerCloseIcon,
  commentsUlEl,
  deleteStoryBtn,
  editorReadOnlyBody,
  editorReadOnlyHeader,
  filterBlack,
  followBtn,
  followingBtn,
  likedContainer,
  preparedCommentBox,
  readTimeSpan,
  seeCommentDiv,
} from "../elements/read-story";
import { readPageFollowBtnClick, readPageFollowingBtnClick } from "../follow";
import { getIdParam } from "../helper";

const getFomattedCommentDate = (createdAt: any): string => {
  const createdAtNum = +createdAt;
  if (!isNaN(createdAtNum)) {
    const inMilliSec = createdAtNum * 1000;
    const unixNow = new Date().getTime();
    const diff = unixNow - inMilliSec;
    const diffMin = Math.ceil(diff / 60 / 1000);
    if (diffMin === 1) {
      return `just now`;
    }
    if (diffMin < 60) {
      return `${diffMin} minutes ago`;
    }
    return new Date(inMilliSec).toLocaleString();
  } else {
    return "unknown";
  }
};

export const drawNewComment = (comment: any) => {
  if (commentsUlEl) {
    const {
      commentId,
      userId,
      avatarUrl,
      createdAt,
      text,
      username,
      isAuthorized,
    } = comment;
    const commentDate = getFomattedCommentDate(createdAt);
    const textEl = document.createElement("p");
    const timestampEl = document.createElement("div");
    const usernameEl = document.createElement("a");
    const infoEl = document.createElement("div");
    const avatarImgEl = document.createElement("img");
    const avatarFrameEl = document.createElement("div");
    let deleteIcon;
    const headerEl = document.createElement("header");
    const liEl = document.createElement("li");
    textEl.className = "comment__text";
    timestampEl.className = "comment__timestamp";
    usernameEl.className = "comment__creator-name";
    infoEl.className = "comment__info _flex-c-sb";
    avatarImgEl.className = "_avatar-img";
    avatarFrameEl.className = "_avatar-frame";
    headerEl.className = "comment__header _flex-cs";
    liEl.className = "drawer-comment__comment-container";
    liEl.id = commentId;
    textEl.innerText = text;
    timestampEl.innerText = commentDate;
    usernameEl.innerText = username;
    usernameEl.href = `/user-home/${userId}`;
    avatarImgEl.src = avatarUrl;
    infoEl.append(usernameEl);
    infoEl.append(timestampEl);
    avatarFrameEl.append(avatarImgEl);
    if (isAuthorized) {
      deleteIcon = document.createElement("i");
      deleteIcon.className = "far fa-trash-alt comment__delete-icon";
      deleteIcon.addEventListener("click", () =>
        deleteComment(liEl, commentId)
      );
      headerEl.append(deleteIcon);
    }
    headerEl.append(avatarFrameEl);
    headerEl.append(infoEl);
    liEl.append(headerEl);
    liEl.append(textEl);
    commentsUlEl.prepend(liEl);
  }
};

export const clearCommentBox = () => {
  preparedCommentBox && (preparedCommentBox.innerHTML = "");
};

const clearComments = () => {
  commentsUlEl && (commentsUlEl.innerHTML = "");
};

const computeAndPasteReadTime = (blocks: OutputBlockData[]) => {
  if (!readTimeSpan) {
    return;
  }
  let wordCount = 0;
  blocks.forEach((block) => {
    if (block.type === "paragraph") {
      const words: string[] = block.data.text.split(" ");
      wordCount += words.length;
    }
    if (block.type === "code") {
      let words: string[] = block.data.code?.split(" ");
      words = words?.filter(
        (word) =>
          word !== "" &&
          !word.includes("=") &&
          word != "()" &&
          word != "(" &&
          word != ")" &&
          word != "{" &&
          word != "}" &&
          word != "<" &&
          word != ">"
      );
      wordCount += words?.length;
    }
  });
  const readTimeMinute = Math.ceil(wordCount / 200);
  !isNaN(readTimeMinute) &&
    (readTimeSpan.innerText = `${readTimeMinute} min read`);
};

const overrideEditorJsStyleHeader = () => {
  const x = editorReadOnlyHeader?.querySelector(
    ".codex-editor__redactor"
  ) as HTMLElement | null;
  x!.style.paddingBottom = "1rem";
};

export const overrideEditorJsStyleBody = () => {
  const quotes = editorReadOnlyBody?.querySelectorAll(".cdx-quote__text") as
    | NodeListOf<HTMLElement>
    | undefined;
  quotes?.forEach((quote) => {
    quote.style.minHeight = "0px";
  });
};

const initEditorReadOnly = async (storyId: string) => {
  try {
    const { data: blocks } = await Axios.get(`/api/blocks/${storyId}`);
    const header = blocks.shift();
    const headerEditor = new EditorJS({
      //   readOnly: true, // This occurs error on 2.19.0 version. It's on github issue https://github.com/codex-team/editor.js/issues/1400
      ...EDITORJS_CONFIG,
      holder: "editor-readOnly__header",
      data: { blocks: [header] },
    });
    const bodyEditor = new EditorJS({
      //   readOnly: true, // This occurs error on 2.19.0 version. It's on github issue https://github.com/codex-team/editor.js/issues/1400
      ...EDITORJS_CONFIG,
      holder: "editor-readOnly__body",
      data: { blocks },
    });
    headerEditor.isReady.then(async () => {
      await headerEditor.readOnly.toggle(true);
      overrideEditorJsStyleHeader();
    });
    bodyEditor.isReady.then(async () => {
      await bodyEditor.readOnly.toggle(true);
      computeAndPasteReadTime(blocks);
      overrideEditorJsStyleBody();
    });
  } catch {
    alert("Failed to initialize editor. Please try again.");
  }
};

const closeCommentDrawer = () => {
  commentDrawer && (commentDrawer.style.right = "-26rem");
  filterBlack && (filterBlack.style.backgroundColor = "transparent");
  filterBlack && (filterBlack.style.pointerEvents = "none");
  clearCommentBox();
  clearComments();
};

const openCommentDrawer = () => {
  commentDrawer && (commentDrawer.style.right = "0rem");
  filterBlack && (filterBlack.style.backgroundColor = "black");
  filterBlack && (filterBlack.style.pointerEvents = "all");
};

const getComments = async () => {
  const storyId = getIdParam("read-story");
  try {
    const { status, data } = await Axios.get(`/api/comment/${storyId}`);
    if (status < 300) {
      data.forEach((comment: any) => {
        drawNewComment(comment);
      });
    }
  } catch {
    alert("Failed to load comments. Please try again.");
  }
};

const toggleLike = async () => {
  const storyId = getIdParam("read-story");
  const childIcon = likedContainer?.querySelector("i");
  const childSpan = likedContainer?.querySelector("span");
  if (!likedContainer || !childIcon || !childSpan) {
    return;
  }
  const likedCount = parseInt(childSpan.innerText.replace(/\,/, ""));
  if (isNaN(likedCount)) {
    console.error("wrong like count format");
    return;
  }
  try {
    likedContainer.style.pointerEvents = "none";
    const { status, data } = await Axios.patch(`/api/toggle-like/${storyId}`);
    if (status < 300) {
      const result = +data;
      if (isNaN(result)) {
        return;
      }
      if (result > 0) {
        childSpan.innerText = (likedCount + 1).toLocaleString();
        childIcon.className = childIcon.className.replace("far", "fas");
      } else if (result < 0) {
        childSpan.innerText = (likedCount - 1).toLocaleString();
        childIcon.className = childIcon.className.replace("fas", "far");
      } else {
        return;
      }
    }
  } catch {
    alert("Failed to like/unlike the story. Please try again.");
  } finally {
    likedContainer.style.pointerEvents = "auto";
  }
};

const handleBookmark = async () => {
  const storyId = getIdParam("read-story");
  const childIcon = bookmarkContainer?.querySelector("i");
  if (childIcon) {
    if (childIcon.className.includes("false")) {
      try {
        const { status } = await Axios.patch(`/api/bookmark/${storyId}`);
        if (status < 300) {
          childIcon.className = childIcon.className
            .replace("false", "true")
            .replace("far", "fas");
        }
      } catch {
        alert("Failed to bookmark. Please try again.");
      }
    } else if (childIcon.className.includes("true")) {
      try {
        const { status } = await Axios.patch(`/api/disbookmark/${storyId}`);
        if (status < 300) {
          childIcon.className = childIcon.className
            .replace("true", "false")
            .replace("fas", "far");
        }
      } catch {
        alert("Failed to disbookmark. Please try again.");
      }
    }
  }
};

const init = async () => {
  const storyId = getIdParam("read-story");
  await initEditorReadOnly(storyId);

  filterBlack?.addEventListener("click", closeCommentDrawer);
  seeCommentDiv?.addEventListener("click", openCommentDrawer);
  seeCommentDiv?.addEventListener("click", getComments);
  commentDrawerCloseIcon?.addEventListener("click", closeCommentDrawer);
  likedContainer?.addEventListener("click", toggleLike);
  bookmarkContainer?.addEventListener("click", handleBookmark);
  followBtn?.addEventListener("click", readPageFollowBtnClick);
  followingBtn?.addEventListener("click", readPageFollowingBtnClick);
  deleteStoryBtn?.addEventListener("click", deleteStory);
};

if (document.location.pathname.includes("read-story")) {
  init();
}
